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
)

type gantry struct {
	sprite *Sprite
	level  *Level
	startx int
}

func (o *gantry) HandleEvent(ev tcell.Event) bool {
	s := o.sprite
	switch ev := ev.(type) {
	case *EventSpriteMove:
		x, _ := o.sprite.Position()
		if o.startx != 0 && o.startx-5 >= x {
			o.sprite.SetVelocity(0, 0)
		}
	case *EventCollision:
		switch ev.Collider().Layer() {
		case LayerShot, LayerPlayer:
			o.sprite.Hide()
			x, y, _, _ := s.Bounds()
			props := GameObjectProps{}
			props.PropSetInt("x", x)
			props.PropSetInt("y", y)
			props.PropSetInt("count", 2)
			MakeGameObject(o.level, "explosion", props)
		}
	case *EventLevelStart:
		o.sprite.SetVelocity(-3.0, 0)
		o.sprite.SetFrame("RETRACT")
		o.startx, _ = o.sprite.Position()
	}
	return false
}

func makeGantry(level *Level, props GameObjectProps) error {

	x := props.PropInt("x", 0)
	y := props.PropInt("y", 0)
	frame := props.PropString("frame", "START")
	sname := props.PropString("sprite", "Gantry")

	sprite := GetSprite(sname)
	o := &gantry{sprite: sprite, level: level}
	sprite.SetLayer(LayerHazard)
	sprite.Watch(o)
	sprite.SetFrame(frame)
	sprite.SetPosition(x, y)
	level.AddSprite(sprite)
	return nil
}

func init() {
	RegisterGameObjectMaker("gantry", makeGantry)
}
