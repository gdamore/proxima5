// Copyright 2016 Garrett D'Amore
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"time"

	"github.com/gdamore/tcell"
)

type ship struct {
	ship     *Sprite
	xthrust  *Sprite
	ythrust  *Sprite
	level    *Level
	lastx    int
	lasty    int
	ximpulse float64
	yimpulse float64
	bspeed   float64
	lastgrav time.Time
	launched bool
	dead     bool
}

func (o *ship) shoot() {
	// TBD: shot limiter
	x1, y, x2, _ := o.ship.Bounds()
	_, vy := o.ship.Velocity()
	props := GameObjectProps{}
	props.PropSetInt("x", x1+(x2-x1)/2)
	props.PropSetInt("y", y-1)
	props.PropSetFloat64("vy", vy-o.bspeed)
	props.PropSetInt("delay", 10)
	MakeGameObject(o.level, "bullet", props)
}

func (o *ship) thrustUp() {
	vx, vy := o.ship.Velocity()
	vy -= o.yimpulse
	o.ship.SetVelocity(vx, vy)
	o.ythrust.SetFrame("UP")
	if !o.launched {
		if sprite := GetSprite("Exhaust"); sprite != nil {
			x, y := o.ship.Position()
			sprite.SetFrame("LIFTOFF")
			sprite.SetPosition(x, y)
			sprite.SetLayer(LayerThrust)
			o.level.AddSprite(sprite)
		}
		o.launched = true
	}
}

func (o *ship) thrustLeft() {
	vx, vy := o.ship.Velocity()
	vx -= o.ximpulse
	o.ship.SetVelocity(vx, vy)
	o.xthrust.SetFrame("LEFT")
}

func (o *ship) thrustRight() {
	vx, vy := o.ship.Velocity()
	vx += o.ximpulse
	o.ship.SetVelocity(vx, vy)
	o.xthrust.SetFrame("RIGHT")
}

func (o *ship) thrustDown() {
	vx, vy := o.ship.Velocity()
	vy += o.yimpulse
	o.ship.SetVelocity(vx, vy)
	o.ythrust.SetFrame("DOWN")
}

func (o *ship) destroy() {
	o.dead = true
	o.ship.Hide()

	x, y, _, _ := o.ship.Bounds()
	props := GameObjectProps{}
	props.PropSetInt("x", x)
	props.PropSetInt("y", y)
	props.PropSetInt("count", 3)
	MakeGameObject(o.level, "explosion", props)
	props = GameObjectProps{}
	props.PropSetInt("x", x)
	props.PropSetInt("y", y)
	props.PropSetInt("count", 3)
	MakeGameObject(o.level, "smexplosion", props)
	o.level.RemoveSprite(o.ship)
	o.level.HandleEvent(&EventPlayerDeath{})
}

func (o *ship) adjustView() {
	x1, y1, x2, y2 := o.ship.Bounds()
	o.level.MakeVisible(x2+20, y2+8)
	o.level.MakeVisible(x1-20, y1-12)
	if o.ythrust != nil {
		o.ythrust.SetPosition(x1+(x2-x1)/2, y1)
	}
	if o.xthrust != nil {
		o.xthrust.SetPosition(x1+(x2-x1)/2, y1)
	}
}

func (o *ship) HandleEvent(ev tcell.Event) bool {
	if o.dead {
		return false
	}
	switch ev := ev.(type) {
	case *EventSpriteAccelerate:
		if ev.s != o.ship {
			return false
		}
		vx, _ := o.ship.Velocity()
		if vx >= 1.0 {
			o.ship.SetFrame("RIGHT")
		} else if vx <= -1.0 {
			o.ship.SetFrame("LEFT")
		} else {
			o.ship.SetFrame("FWD")
		}
	case *EventSpriteMove:
		// We don't let ship leave the map
		x, y := o.ship.Position()
		ox, oy := x, y
		vx, vy := o.ship.Velocity()
		w, h := o.level.Size()
		if x < 0 {
			x = 0
			if vx < 0 {
				vx = 0
			}
		} else if x >= w {
			x = w - 1
			if vx > 0 {
				vx = 0
			}
		}
		if y < 0 {
			y = 0
			if vy < 0 {
				vy = 0
			}
		} else if y >= h {
			y = h - 1
			if vy > 0 {
				vy = 0
			}
		}
		if ox != x || oy != y {
			o.ship.SetPosition(x, y)
			o.ship.SetVelocity(vx, vy)
		}

		if y == 0 {
			o.dead = true
			o.level.HandleEvent(&EventLevelComplete{})
		}
		o.adjustView()
	case *EventGravity:
		now := ev.When()
		if !o.lastgrav.IsZero() {
			vx, vy := o.ship.Velocity()
			frac := float64(now.Sub(o.lastgrav))
			frac /= float64(time.Second)
			vy += ev.Accel() * frac
			o.ship.SetVelocity(vx, vy)
		}
		o.lastgrav = now

	case *EventCollision:
		switch ev.Collider().Layer() {
		case LayerTerrain, LayerHazard, LayerShot:
			o.destroy()
		case LayerPad:
			// if we're on the pad, and not too
			// fast, then stay on the pad.
			// TODO: probably the max velocity (4.0)
			// should be tunable.
			vx, vy := o.ship.Velocity()
			x, y := o.ship.Position()
			if vx == 0 && vy > 0 && vy < 4.0 {
				y--
				vy = 0
				o.ship.SetPosition(x, y)
				o.ship.SetVelocity(vx, vy)
				o.launched = false
			} else {
				o.destroy()
			}
		}

	case *EventTimesUp:
		o.destroy()

	case *tcell.EventKey:
		switch ev.Key() {

		case tcell.KeyLeft:
			o.thrustLeft()
			return true

		case tcell.KeyRight:
			o.thrustRight()
			return true

		case tcell.KeyUp:
			o.thrustUp()
			return true

		case tcell.KeyDown:
			o.thrustDown()
			return true

		case tcell.KeyRune:
			switch ev.Rune() {
			case ' ':
				o.shoot()
				return true
			case 'j', 'J':
				o.thrustLeft()
				return true
			case 'k', 'K':
				o.thrustRight()
				return true
			case 'i', 'I':
				o.thrustUp()
				return true
			case 'm', 'M':
				o.thrustDown()
				return true
			}
		}
	case *tcell.EventResize:
		x, y := o.ship.Position()
		o.level.Center(x, y)
		o.adjustView()
	}
	return false
}

func makeShip(level *Level, props GameObjectProps) error {

	x := props.PropInt("x", 0)
	y := props.PropInt("y", 0)
	velx := props.PropFloat64("vx", 0)
	vely := props.PropFloat64("vy", 0)
	frame := props.PropString("frame", "FWD")
	sname := props.PropString("sprite", "Ship")
	xthrust := props.PropString("xthrust", "Thrust")
	ythrust := props.PropString("ythrust", "Thrust")

	sprite := GetSprite(sname)
	o := &ship{ship: sprite, level: level}

	// How much impulse does each engine fire give us?
	o.yimpulse = props.PropFloat64("yimpulse", 3.5)
	o.ximpulse = props.PropFloat64("ximpulse", 4.5)
	o.bspeed = props.PropFloat64("vshot", 20.0)

	sprite.SetLayer(LayerPlayer)
	sprite.Watch(o)
	sprite.SetFrame(frame)
	sprite.SetPosition(x, y)
	sprite.SetVelocity(velx, vely)
	level.AddSprite(sprite)

	o.xthrust = GetSprite(xthrust)
	o.ythrust = GetSprite(ythrust)
	level.AddSprite(o.xthrust)
	level.AddSprite(o.ythrust)
	o.xthrust.SetLayer(LayerThrust)
	o.ythrust.SetLayer(LayerThrust)
	o.level.Center(x, y)
	o.adjustView()

	return nil
}

func init() {
	RegisterGameObjectMaker("ship", makeShip)
}
