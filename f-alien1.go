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

package proxima5

import (
	"github.com/gdamore/tcell"
)

type alien1 struct {
	sprite *Sprite
	level  *Level
}

func (o *alien1) HandleEvent(ev tcell.Event) bool {
	s := o.sprite
	switch ev := ev.(type) {
	case *EventSpriteAccelerate:
		if ev.s != s {
			return false
		}
		vx, _ := s.Velocity()
		if vx > 0 {
			s.SetFrame("F1")
		} else {
			s.SetFrame("R1")
		}
	case *EventCollision:
		switch ev.Collider().Layer() {
		case LayerTerrain, LayerHazard:
			x, y, _, _ := s.Bounds()
			vx, vy := s.Velocity()
			vx = -vx
			vy = -vy
			s.SetVelocity(vx, vy)
			if vx < 0 {
				x--
			} else if vx > 0 {
				x++
			}
			if vy < 0 {
				y--
			} else if vy > 0 {
				y++
			}
			s.SetPosition(x, y)
		case LayerShot, LayerPlayer:
			s.Hide()
			x, y, _, _ := s.Bounds()
			props := GameObjectProps{}
			props.PropSetInt("x", x)
			props.PropSetInt("y", y)
			props.PropSetInt("count", 1)
			MakeGameObject(o.level, "smexplosion", props)
			o.level.RemoveSprite(o.sprite)
		}
	}
	return false
}

func makeAlien1(level *Level, props GameObjectProps) error {

	x := props.PropInt("x", 0)
	y := props.PropInt("y", 0)
	velx := props.PropFloat64("vx", 3.5)
	vely := props.PropFloat64("vy", 0)
	frame := props.PropString("frame", "F1")
	sname := props.PropString("sprite", "Alien1")

	sprite := GetSprite(sname)
	o := &alien1{sprite: sprite, level: level}
	sprite.SetLayer(LayerHazard)
	sprite.Watch(o)
	sprite.SetFrame(frame)
	sprite.SetPosition(x, y)
	sprite.SetVelocity(velx, vely)
	level.AddSprite(sprite)
	return nil
}

func init() {
	RegisterGameObjectMaker("alien1", makeAlien1)
}
