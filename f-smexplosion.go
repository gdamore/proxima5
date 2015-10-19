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
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

//
// "smexplosion" understands these properties:
//
//  count  - number of little explosion clouds
//  interval - intervalin msec between clouds
//  xdistance - max x distance in cells between explosion clouds
//  ydistance - max y distance in cells between explosion clouds
//
type smExplosion struct {
	level *Level
}

func (o *smExplosion) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *EventSpriteFrame:
		if ev.Frame() == "" {
			o.level.RemoveSprite(ev.Sprite())
		}
	}
	return false
}

func makeSmExplosion(level *Level, props GameObjectProps) error {

	x := props.PropInt("x", 0)
	y := props.PropInt("y", 0)
	o := &smExplosion{level: level}

	frame := props.PropString("frame", "0")
	sname := props.PropString("sprite", "SmallExplosion")
	count := props.PropInt("count", 1)
	interval := props.PropInt("interval", 80)
	xdist := props.PropInt("xdistance", 3)
	ydist := props.PropInt("ydistance", 2)
	delay := props.PropInt("delay", 0)

	lx, ly := x, y
	d := time.Duration(delay) * time.Millisecond
	for i := 0; i < count; i++ {
		sprite := GetSprite(sname)
		if i > 0 {
			d += time.Millisecond * time.Duration(interval)
			lx += rand.Intn(xdist*2) - xdist
			ly += rand.Intn(xdist*2) - ydist
		}
		if i > 0 || delay != 0 {
			sprite.SetFrame("")
			sprite.ScheduleFrame(frame, time.Now().Add(d))
		} else {
			sprite.SetFrame("0")
		}
		sprite.SetPosition(lx, ly)
		sprite.SetLayer(LayerExplosion)
		sprite.Watch(o)
		level.AddSprite(sprite)
	}

	return nil
}

func init() {
	RegisterGameObjectMaker("smexplosion", makeSmExplosion)
}
