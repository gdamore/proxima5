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
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

//
// "explosion" understands these class-specific properties:
//
//  count  - number of little explosion clouds
//  interval - intervalin msec between clouds
//  xdistance - max x distance in cells between explosion clouds
//  ydistance - max y distance in cells between explosion clouds
//
type explosion struct {
	level *Level
}

func (o *explosion) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *EventSpriteFrame:
		if ev.Frame() == "" {
			o.level.RemoveSprite(ev.Sprite())
		}
	}
	return false
}

func makeExplosion(level *Level, props GameObjectProps) error {

	x := props.PropInt("x", 0)
	y := props.PropInt("y", 0)
	o := &explosion{level: level}

	frame := props.PropString("frame", "0")
	sname := props.PropString("sprite", "Explosion")
	count := props.PropInt("count", 1)
	interval := props.PropInt("interval", 80)
	xdist := props.PropInt("xdistance", 5)
	ydist := props.PropInt("ydistance", 3)

	lx, ly := x, y
	d := time.Duration(0)
	for i := 0; i < count; i++ {
		sprite := GetSprite(sname)
		if i > 0 {
			lx += rand.Intn(xdist*2) - xdist
			ly += rand.Intn(xdist*2) - ydist
			d += time.Millisecond * time.Duration(interval)
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
	RegisterGameObjectMaker("explosion", makeExplosion)
}
