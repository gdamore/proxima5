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
	"time"

	"github.com/gdamore/tcell"
)

type hfence struct {
	remit *Sprite
	lemit *Sprite
	beam  *Sprite
	level *Level
	dead  bool
}

func (o *hfence) destroy(primary, alternate *Sprite) {
	o.dead = true
	o.beam.SetFrame("")
	o.remit.SetFrame("RX")
	o.lemit.SetFrame("LX")

	x1, y1, x2, y2 := primary.Bounds()
	props := GameObjectProps{}
	props.PropSetInt("x", x1+(x2-x1)/2)
	props.PropSetInt("y", y1+(y2-y1)/2)
	props.PropSetInt("count", 2)
	MakeGameObject(o.level, "smexplosion", props)

	x1, y1, x2, y2 = alternate.Bounds()
	props = GameObjectProps{}
	props.PropSetInt("x", x1+(x2-x1)/2)
	props.PropSetInt("y", y1+(y2-y1)/2)
	props.PropSetInt("delay", 500)
	props.PropSetInt("count", 1)
	MakeGameObject(o.level, "smexplosion", props)
}

func (o *hfence) HandleEvent(ev tcell.Event) bool {
	if o.dead {
		return false
	}
	switch ev := ev.(type) {
	case *EventCollision:
		switch ev.Collider().Layer() {
		case LayerPlayer:
			switch ev.Target() {
			case o.lemit:
				o.destroy(o.lemit, o.remit)
			case o.remit:
				o.destroy(o.remit, o.lemit)
			}
		case LayerShot:
			switch ev.Target() {
			case o.lemit:
				o.destroy(o.lemit, o.remit)
			case o.remit:
				o.destroy(o.remit, o.lemit)
			case o.beam:
				x, y, _, _ := ev.Collider().Bounds()
				props := GameObjectProps{}
				props.PropSetInt("x", x)
				props.PropSetInt("y", y)
				props.PropSetInt("count", 1)
				props.PropSetString("sprite", "TinyExplosion")
				MakeGameObject(o.level, "explosion", props)
			}
		}
	}
	return false
}

func makeHFence(level *Level, props GameObjectProps) error {

	x := props.PropInt("x", 0)
	y := props.PropInt("y", 0)
	width := props.PropInt("width", 6)
	height := props.PropInt("height", 1)
	offtime := props.PropInt("offtime", 4000)
	ontime := props.PropInt("ontime", 2000)

	if width < 5 || height < 1 {
		// need at least 5 places
		return nil
	}
	sname := props.PropString("sprite", "HFence")

	o := &hfence{level: level}
	o.remit = GetSprite(sname)
	o.lemit = GetSprite(sname)
	o.beam = GetSprite(sname)

	o.remit.Resize(2, height)
	o.lemit.Resize(2, height)
	o.beam.Resize(width-4, height)

	o.lemit.SetFrame("LL")
	o.remit.SetFrame("RR")
	o.beam.SetFrame("OFF")

	o.beam.SetNextFrame("OFF",
		time.Millisecond*time.Duration(offtime), "ON")

	if ontime > 1000 {
		ontime -= 1000
		o.beam.SetNextFrame("ON", 500*time.Millisecond, "ON1")
		o.beam.SetNextFrame("ON1",
			time.Duration(ontime)*time.Millisecond, "ON2")
		o.beam.SetNextFrame("ON2", 500*time.Millisecond, "OFF")
	} else {
		// Too short to reach full charge
		o.beam.SetNextFrame("ON",
			time.Duration(ontime)*time.Millisecond, "OFF")
	}

	o.remit.SetLayer(LayerHazard)
	o.lemit.SetLayer(LayerHazard)
	o.beam.SetLayer(LayerHazard)

	o.lemit.SetPosition(x, y)
	o.remit.SetPosition(x+width-2, y)
	o.beam.SetPosition(x+2, y)

	o.lemit.Watch(o)
	o.remit.Watch(o)
	o.beam.Watch(o)
	o.level = level
	level.AddSprite(o.lemit)
	level.AddSprite(o.remit)
	level.AddSprite(o.beam)

	return nil
}

func init() {
	RegisterGameObjectMaker("hfence", makeHFence)
}
