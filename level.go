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
	"github.com/gdamore/tcell/views"
)

type Level struct {
	name     string
	title    string
	width    int
	height   int
	startx   int
	starty   int
	game     *Game
	layer    int
	maxlayer int
	view     *views.ViewPort
	maxtime  time.Duration
	begin    time.Time
	expired  bool
	gravity  float64
	clock    AlarmClock
	stopped  bool
	manager  *SpriteManager
	terrain  *Sprite
	data     *LevelData
	win      bool
	started  bool
}

func NewLevelFromData(d *LevelData) *Level {
	lvl := &Level{data: d}

	lvl.name = d.Name
	lvl.title = d.Properties.PropString("title", lvl.name)

	timer := d.Properties.PropInt("timer", 300)
	lvl.gravity = d.Properties.PropFloat64("gravity", 0.0)
	lvl.width = d.Properties.PropInt("width", 0)
	lvl.height = d.Properties.PropInt("height", 0)
	if lvl.width == 0 || lvl.height == 0 {
		panic("level has no dimensions!")
	}
	lvl.maxtime = time.Duration(timer) * time.Second
	lvl.manager = NewSpriteManager(lvl.width, lvl.height)
	lvl.clock = NewAlarmClock()

	return lvl
}

func (l *Level) Name() string {
	return l.name
}

func (l *Level) Title() string {
	return l.title
}

func (l *Level) randomViewXY() (int, int) {
	x1, y1, x2, y2 := l.view.GetVisible()
	return rand.Intn(x2-x1+1) + x1, rand.Intn(y2-y1+1) + y1
}

func (l *Level) Update(now time.Time) {

	l.clock.Tick(now)
	l.manager.Update(now)

	if l.gravity != 0 {
		l.HandleEvent(&EventGravity{when: now, accel: l.gravity})
	}

	if l.GetTimer() == 0 && !l.expired && !l.win && l.started {
		l.HandleEvent(&EventTimesUp{})
		l.expired = true
	}
}

func (l *Level) SetView(v *views.ViewPort) {
	l.view = v
	v.SetContentSize(l.width, l.height, true)
	l.manager.SetView(v)
}

func (l *Level) Draw(v *views.ViewPort) {

	l.manager.Draw(v)
}

func (l *Level) Reset() {

	l.manager.Reset()
	l.clock.Reset()

	for name, plist := range l.data.Objects {
		props := GameObjectProps{}
		props["label"] = name
		for k, v := range plist {
			props[k] = v
		}
		MakeGameObject(l, props["class"], props)
	}
	l.begin = time.Now()
	l.stopped = false
	l.started = false
	l.win = false
	l.expired = false
}

func (l *Level) GetTimer() time.Duration {
	if l.win || !l.started {
		return l.maxtime
	}
	d := time.Now().Sub(l.begin)
	d = l.maxtime - d
	if d < 0 {
		d = 0
	}
	return d
}

func (l *Level) SetGame(g *Game) {
	l.game = g
}

func (l *Level) Layer() int {
	return LayerTerrain
}

func (l *Level) Size() (int, int) {
	return l.width, l.height
}

func (l *Level) MakeVisible(x, y int) {
	if l.view != nil {
		l.view.MakeVisible(x, y)
	}
}

func (l *Level) Center(x, y int) {
	if l.view != nil {
		l.view.Center(x, y)
	}
}

func (l *Level) Start() {
	l.started = true
	l.begin = time.Now()
	l.HandleEvent(&EventLevelStart{})
}

func (l *Level) ShowPress() {
	x1, y1, x2, y2 := l.view.GetVisible()
	sprite := GetSprite("PressSpace")
	sprite.SetLayer(LayerDialog)
	sprite.SetFrame("F0")
	sprite.SetPosition(x1+(x2-x1)/2, y1+(y2-y1)/2+2)
	l.manager.AddSprite(sprite)
}

func (l *Level) ShowComplete() {
	x1, y1, x2, y2 := l.view.GetVisible()
	sprite := GetSprite("LevelComplete")
	sprite.SetLayer(LayerDialog)
	sprite.SetFrame("F0")
	sprite.SetPosition(x1+(x2-x1)/2, y1+(y2-y1)/2)
	l.manager.AddSprite(sprite)
}

func (l *Level) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlR:
			// secret reset button
			l.Reset()
		}
	case *tcell.EventMouse:

		offx, offy, _, _ := l.view.GetVisible()
		px1, py1, px2, py2 := l.view.GetPhysical()
		mx, my := ev.Position()
		if mx < px1 || mx > px2 || my < py1 || my > py2 {
			// outside our view
			return false
		}
		if ev.Buttons()&tcell.Button1 != 0 {
			l.stopped = true
			l.view.Center(offx+mx-px1, offy+my-py1)
		} else if ev.Buttons()&tcell.Button2 != 0 {
			l.Reset()
		} else if ev.Buttons()&tcell.Button3 != 0 {
			l.Reset()
		}
		return true

	case *EventPlayerDeath:
		l.game.HandleEvent(ev)
		l.ShowPress()
		return true

	case *EventLevelComplete:
		l.win = true
		l.game.HandleEvent(ev)
		l.ShowComplete()
		l.ShowPress()
		return true

	case *EventGameOver:
		x1, y1, x2, y2 := l.view.GetVisible()
		sprite := GetSprite("GameOver")
		sprite.SetLayer(LayerDialog)
		sprite.SetFrame("F0")
		sprite.SetPosition(x1+(x2-x1)/2, y1+(y2-y1)/2)
		l.manager.AddSprite(sprite)

	case *EventTimesUp:
		bw := GetSprite("Blastwave")
		bw.Resize(l.width, l.height)
		bw.ScheduleFrame("0", time.Now().Add(2*time.Second))
		l.AddSprite(bw)

		dur := time.Duration(0)
		for i := 0; i < 100; i++ {
			x, y := l.randomViewXY()
			sprite := GetSprite("Explosion")
			sprite.ScheduleFrame("0", time.Now().Add(dur))
			sprite.SetPosition(x, y)
			l.AddSprite(sprite)
			dur += time.Millisecond * 5

			for j := 0; j < 4; j++ {
				x, y = l.randomViewXY()
				sprite = GetSprite("SmallExplosion")
				sprite.SetPosition(x, y)
				sprite.ScheduleFrame("0", time.Now().Add(dur))
				l.AddSprite(sprite)
				dur += time.Millisecond * 50
			}
		}
	}
	if l.stopped {
		return false
	}

	return l.manager.HandleEvent(ev)
}

func (l *Level) AddSprite(sprite *Sprite) {
	l.manager.AddSprite(sprite)
}

func (l *Level) RemoveSprite(sprite *Sprite) {
	l.manager.RemoveSprite(sprite)
}

func (l *Level) AddAlarm(d time.Duration, h EventHandler) *Alarm {
	return l.clock.Schedule(d, h)
}

func (l *Level) RemoveAlarm(a *Alarm) {
	l.clock.Cancel(a)
}
