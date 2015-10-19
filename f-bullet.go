// Copyright 2015 Garrett D'Amore
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
	"github.com/gdamore/tcell"
	"time"
)

type bullet struct {
	sprite *Sprite
	level  *Level
	alarm  *Alarm
}

func (o *bullet) destroy() {
	o.sprite.Hide()
	o.level.RemoveSprite(o.sprite)
	o.level.RemoveAlarm(o.alarm)
}

func (o *bullet) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *EventSpriteAccelerate:
		if ev.s != o.sprite {
			return false
		}

		vx, _ := o.sprite.Velocity()
		if vx > 0 {
			o.sprite.SetFrame("H")
		} else {
			o.sprite.SetFrame("V")
		}
	case *EventSpriteMove:
		if ev.s != o.sprite {
			return false
		}

		x, y, _, _ := o.sprite.Bounds()
		w, h := o.level.Size()
		if x < 0 || y < 0 || x >= w || y >= h {
			o.destroy()
		}
	case *EventCollision:
		switch ev.Collider().Layer() {

		case LayerTerrain, LayerHazard, LayerPlayer, LayerExplosion:
			// Impact with most solid objects removes the shot.  The
			// impacted object is responsible for painting any explosive
			// effect.
			o.destroy()
		}
	case *EventAlarm:
		o.destroy()
	}
	return false
}

func makeBullet(level *Level, props GameObjectProps) error {

	x := props.PropInt("x", 0)
	y := props.PropInt("y", 0)
	velx := props.PropFloat64("vx", 0)
	vely := props.PropFloat64("vy", -20.0)
	frame := props.PropString("frame", "V")
	sname := props.PropString("sprite", "Bullet")
	life := props.PropInt("lifetime", 2000)
	sprite := GetSprite(sname)
	o := &bullet{sprite: sprite, level: level}

	if life > 0 {
		a := level.AddAlarm(time.Millisecond*time.Duration(life), o)
		o.alarm = a
	}

	sprite.SetLayer(LayerShot)
	sprite.Watch(o)
	sprite.SetFrame(frame)
	sprite.SetPosition(x, y)
	sprite.SetVelocity(velx, vely)
	level.AddSprite(sprite)
	return nil
}

func init() {
	RegisterGameObjectMaker("bullet", makeBullet)
}
