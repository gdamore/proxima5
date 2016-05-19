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
	"errors"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"sync"
	"time"
)

// go-proxima (Escape from Proxima Five) is a terminal-oriented game,
// utilizing Unicode characters and text terminal display capabilties,
// demonstrating  the capabilities of tcell.   It needs at leasst an 80x24
// column terminal to run well.

type Game struct {
	quitq    chan struct{}
	screen   tcell.Screen
	eventq   chan tcell.Event
	errmsg   string
	quitone  sync.Once
	level    *Level
	lview    *views.ViewPort
	sview    *views.ViewPort
	sbar     *views.TextBar
	lives    int
	gameover bool
	started  bool

	sync.Mutex
}

func (g *Game) Init() error {
	g.lives = 5
	if screen, err := tcell.NewScreen(); err != nil {
		return err
	} else if err = screen.Init(); err != nil {
		return err
	} else {
		screen.SetStyle(tcell.StyleDefault.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorWhite))
		g.screen = screen
	}

	// XXX: Add a main screen
	g.screen.EnableMouse()
	g.level = GetLevel("level1")
	if g.level == nil {
		g.screen.Fini()
		return errors.New("Cannot find data (did you run rebuild.sh?)")
	}
	g.lview = views.NewViewPort(g.screen, 0, 1, -1, -1)
	g.level.SetView(g.lview)
	g.level.SetGame(g)

	g.sview = views.NewViewPort(g.screen, 0, 0, -1, 1)
	g.sbar = views.NewTextBar()
	g.sbar.SetView(g.sview)

	g.quitq = make(chan struct{})
	g.eventq = make(chan tcell.Event)

	g.level.Reset()
	g.level.ShowPress()

	RegisterFallbacks(g.screen)

	return nil
}

func (g *Game) Quit() {
	g.quitone.Do(func() {
		close(g.quitq)
	})
}

func (g *Game) Draw() {
	g.Lock()
	sbnorm := tcell.StyleDefault.
		Background(tcell.ColorAqua).
		Foreground(tcell.ColorBlack)
	sbwarn := tcell.StyleDefault.
		Background(tcell.ColorAqua).
		Foreground(tcell.ColorRed)
	sbalert := tcell.StyleDefault.
		Background(tcell.ColorRed).
		Foreground(tcell.ColorWhite)

	g.sbar.SetStyle(sbnorm)
	timer := g.level.GetTimer()
	times := fmt.Sprintf("%02d:%02d.%1d",
		timer/time.Minute,
		(timer%time.Minute)/time.Second,
		timer%time.Second/(100*time.Millisecond))

	g.sbar.SetCenter(g.level.Title(), sbnorm)
	if timer < 30*time.Second {
		g.sbar.SetLeft(times, sbwarn)
		if timer < 10*time.Second {
			g.sbar.SetCenter(
				fmt.Sprintf("-=: BLASTWAVE IN %s :=-",
					times), sbalert)
		}
	} else {
		g.sbar.SetLeft(times, sbnorm)
	}

	if g.gameover {
		g.sbar.SetCenter("-=: G A M E  O V E R :=- ", sbwarn)
	}

	right := ""
	if g.lives > 5 {
		right += fmt.Sprintf("%d x ♥", g.lives)
	} else {
		for i := 0; i < g.lives; i++ {
			right += " ♥"
		}
	}
	right += " "

	g.sbar.SetRight(right, sbwarn)

	g.level.Draw(g.lview)
	g.sbar.Draw()
	g.screen.Show()
	g.Unlock()
}

func (g *Game) Run() error {

	go g.EventPoller()
	go g.Updater()
loop:
	for {
		g.Draw()
		select {
		case <-g.quitq:
			break loop
		case <-time.After(time.Millisecond * 10):
		case ev := <-g.eventq:
			g.HandleEvent(ev)
		}
	}

	// Inject a wakeup interrupt
	iev := tcell.NewEventInterrupt(nil)
	g.screen.PostEvent(iev)

	g.screen.Fini()
	// wait for updaters to finish
	if g.errmsg != "" {
		return errors.New(g.errmsg)
	}
	return nil
}

func (g *Game) Error(msg string) {
	g.errmsg = msg
	g.Quit()
}

func (g *Game) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		g.lview.Resize(0, 1, -1, -1)
		g.sview.Resize(0, 0, -1, 1)
		g.level.HandleEvent(ev)
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			g.Quit()
			return true
		}
		if !g.started {
			if ev.Key() == tcell.KeyEnter {
				if g.gameover {
					g.lives = 5
				}
				g.level.Reset()
				g.level.Start()
				g.started = true
				g.gameover = false
			}
			// eat all keys until level starts
			return true
		}
	case *EventPlayerDeath:
		g.lives--
		g.started = false
		if g.lives == 0 {
			g.gameover = true
			g.level.HandleEvent(&EventGameOver{})
		}
		return true

	case *EventLevelComplete:
		g.lives++ // bonus life (for now)
		g.started = false
		return true
	}
	if !g.level.HandleEvent(ev) {
		return true
	}
	return true
}

func (g *Game) Updater() {
	for {
		select {
		case <-g.quitq:
			return
		case <-time.After(time.Millisecond * 10):
			g.Lock()
			g.level.Update(time.Now())
			g.Unlock()
		}
	}
}

func (g *Game) EventPoller() {
	for {
		select {
		case <-g.quitq:
			return
		default:
		}
		ev := g.screen.PollEvent()
		if ev == nil {
			return
		}
		select {
		case <-g.quitq:
			return
		case g.eventq <- ev:
		}
	}
}
