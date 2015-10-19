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

type spriteFrame struct {
	runes     []rune
	styles    []tcell.Style
	nextFrame string
	timer     time.Duration
}

type Sprite struct {
	originx    int
	originy    int
	posx       int
	posy       int
	fposx      float64
	fposy      float64
	velx       float64
	vely       float64
	ptime      time.Time
	width      int
	height     int
	frame      string
	ftime      time.Time
	layer      int
	handlers   map[EventHandler]struct{}
	sched      time.Time
	schedFrame string

	// allow for multiple frames
	frames map[string]*spriteFrame
}

func NewSprite(w, h int) *Sprite {
	s := &Sprite{}
	// no display!
	s.frame = ""
	s.width = w
	s.height = h
	s.handlers = make(map[EventHandler]struct{})
	s.frames = make(map[string]*spriteFrame)

	// default is to exist at layer 1
	s.layer = 1
	return s
}

func (s *Sprite) SetOrigin(x, y int) {
	// It is possible to set the origin outside of the
	// image.  Don't do that!
	s.originx, s.originy = x, y
}

// AddFrame adds a frame, and return the frame number.
func (s *Sprite) AddFrame(name string, runes []rune, styles []tcell.Style) {
	frame := &spriteFrame{}
	frame.runes = make([]rune, 0, s.width*s.height)
	frame.styles = make([]tcell.Style, 0, s.width*s.height)
	for _, r := range runes {
		frame.runes = append(frame.runes, r)
	}
	for _, st := range styles {
		frame.styles = append(frame.styles, st)
	}
	s.frames[name] = frame
}

// SetNextFrame is used to arrange for automatic frame transitions.
// At the specified duration after this frame is made active, the frame
// will advance to the frame specified by next.  If next is -1, then
// the sprite will be made invisible.  If the duration is is zero, then
// no automatic transition will occur.
func (s *Sprite) SetNextFrame(name string, d time.Duration, next string) {
	if frame, ok := s.frames[name]; ok {
		frame.nextFrame = next
		frame.timer = d
	}
}

func (s *Sprite) SetFrame(frame string) {
	s.frame = frame
	s.ftime = time.Now()
	ev := &EventSpriteFrame{when: s.ftime, frame: frame, s: s}
	s.HandleEvent(ev)
}

// ScheduleFrame arranges to change the frame to the named frame
// at some point near the given time.  The purpose is to allow
// deferring the initial display of a sprite.  Only a single
// frame change may be scheduled this way.
func (s *Sprite) ScheduleFrame(frame string, when time.Time) {
	s.sched = when
	s.schedFrame = frame
}

// UnscheduleFrame cancels the scheduled change.
func (s *Sprite) UnscheduleFrame() {
	s.ScheduleFrame("", time.Time{})
}

func (s *Sprite) Frame() string {
	return s.frame
}

func (s *Sprite) Position() (int, int) {
	return s.posx, s.posy
}

func (s *Sprite) SetPosition(x, y int) {
	s.posx, s.posy = x, y
	s.fposx, s.fposy = float64(x), float64(y)
	s.ptime = time.Now()
	ev := &EventSpriteMove{when: s.ptime, s: s, x: x, y: y}
	s.HandleEvent(ev)
}

// Set the velocity of the sprite.  The units are cells per second.
// This cancels any previous but uncalculated motion, so don't call
// this unless the values are actually different, otherwise the motion
// may appear choppier than it should.
func (s *Sprite) SetVelocity(velx, vely float64) {
	s.velx = velx
	s.vely = vely
	s.ptime = time.Now()
	ev := &EventSpriteAccelerate{when: s.ptime, s: s, x: velx, y: vely}
	s.HandleEvent(ev)
}

// Velocity returns the velocity of the sprite in cells per second.
func (s *Sprite) Velocity() (float64, float64) {
	return s.velx, s.vely
}

func (s *Sprite) Update(now time.Time) {
	if !s.sched.IsZero() && s.sched.Before(now) {
		s.sched = time.Time{}
		s.frame = s.schedFrame
		s.ftime = now
		ev := &EventSpriteFrame{s: s, when: now, frame: s.frame}
		s.HandleEvent(ev)
	} else if frame, ok := s.frames[s.frame]; ok {
		if frame.timer != 0 && now.Sub(s.ftime) > frame.timer {
			s.frame = frame.nextFrame
			s.ftime = now
			ev := &EventSpriteFrame{s: s, when: now, frame: s.frame}
			s.HandleEvent(ev)
		}
	}
	if s.velx != 0 || s.vely != 0 {
		fac := float64(now.Sub(s.ptime)) / float64(time.Second)
		s.ptime = now

		s.fposx += s.velx * fac
		s.fposy += s.vely * fac
		ox, oy := s.posx, s.posy
		s.posx = int(s.fposx)
		s.posy = int(s.fposy)

		if ox != s.posx || oy != s.posy {
			ev := &EventSpriteMove{s: s,
				when: now, x: s.posx, y: s.posy}
			s.HandleEvent(ev)
		}
	}
}

func (s *Sprite) Draw(view View) {

	frame, ok := s.frames[s.frame]
	if !ok {
		return
	}
	styles := frame.styles
	runes := frame.runes
	if runes == nil {
		return
	}
	offx, offy := s.posx-s.originx, s.posy-s.originy
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			i := y*s.width + x
			if len(runes) < i {
				return
			}
			if runes[i] == 0 {
				continue
			}
			style := tcell.StyleDefault
			if len(styles) >= i {
				style = styles[i]
			}
			view.SetContent(x+offx, y+offy,
				runes[i], nil, style)
		}
	}
}

// Resize resizes a sprite.  If the sprite shrinks then the cells
// are clipped from the end.  If the sprite grows, then new cells
// are added by copying to the last to the end.
func (s *Sprite) Resize(w, h int) {
	ow, oh := s.width, s.height
	for _, f := range s.frames {
		nr := make([]rune, w*h)
		ns := make([]tcell.Style, w*h)
		if ow == 0 || oh == 0 {
			f.runes = nr
			f.styles = ns
			continue
		}
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if y >= oh {
					nr[y*w+x] = nr[(y-1)*w+x]
					ns[y*w+x] = ns[(y-1)*w+x]
				} else if x >= ow {
					nr[y*w+x] = nr[y*w+(x-1)]
					ns[y*w+x] = ns[y*w+(x-1)]
				} else {
					nr[y*w+x] = f.runes[y*ow+x]
					ns[y*w+x] = f.styles[y*ow+x]
				}
			}
		}
		f.runes = nr
		f.styles = ns
	}
	s.width = w
	s.height = h

}

// Collides checks if a collision at the given coordinates occurs.
func (s *Sprite) Collides(x, y int) bool {
	x -= s.posx
	x += s.originx
	y -= s.posy
	y += s.originy
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return false
	}
	if frame, ok := s.frames[s.frame]; ok {
		runes := frame.runes
		i := y*s.width + x
		if i >= len(runes) {
			return false
		}
		if runes[i] == 0 {
			return false
		}
		return true
	}
	return false
}

// Bounds returns the bounds of the sprite, for collision
// detection purposes (See Collider interface)
func (s *Sprite) Bounds() (x1, y1, x2, y2 int) {

	x1, y1 = s.posx-s.originx, s.posy-s.originy
	x2, y2 = x1+s.width-1, y1+s.height-1
	return x1, y1, x2, y2
}

// SetLayer sets the layer.  This determines when the object is drawn,
// higher layers are drawn after lower ones.  This is also reported via
// Layer() and may be used in collision event handlers.
func (s *Sprite) SetLayer(layer int) {
	s.layer = layer
}

// Layer returns the layer at which the sprite exists.
func (s *Sprite) Layer() int {
	return s.layer
}

func (s *Sprite) Visible() bool {
	_, ok := s.frames[s.frame]
	return ok
}

func (s *Sprite) Hide() {
	s.SetFrame("")
}

func (s *Sprite) HandleEvent(ev tcell.Event) bool {
	for h := range s.handlers {
		if h.HandleEvent(ev) {
			return true
		}
	}
	return false
}

// Watch adds the handler to the list that will be notified for sprite events.
// Note that handlers may be called in any order.  However, a handler that
// returns true will consume the event, and no other handlers will receive
// it afterwards.  For most events, it probably is best to return false,
// so that other game consumers can see the event as well.
func (s *Sprite) Watch(h EventHandler) {
	s.handlers[h] = struct{}{}
}

// Unwatch removes the handler from the notify list for sprite events.
func (s *Sprite) Unwatch(h EventHandler) {
	delete(s.handlers, h)
}

// In addition to EventCollision, we have these.

// EventSpriteFrame is delivered to the watchers when the frame for
// the sprite changes.  If the frame is -1, then the sprite is no longer
// visible.
type EventSpriteFrame struct {
	when  time.Time
	s     *Sprite
	frame string
}

// When returns the time of the event.
func (ev *EventSpriteFrame) When() time.Time {
	return ev.when
}

// Sprite returns the sprite that changed frames.
func (ev *EventSpriteFrame) Sprite() *Sprite {
	return ev.s
}

// Frame returns the new frame number.  It will be -1 if the
// sprite is no longer visible.  The caller is presumed to know
// the details of the individual frames.
func (ev *EventSpriteFrame) Frame() string {
	return ev.frame
}

// EventSpriteMove is delivered ot the watchers when the sprite
// position changes.
type EventSpriteMove struct {
	when time.Time
	s    *Sprite
	x, y int
}

// When returns the time of the event.
func (ev *EventSpriteMove) When() time.Time {
	return ev.when
}

// Sprite returns the sprite that moved.
func (ev *EventSpriteMove) Sprite() *Sprite {
	return ev.s
}

// Position returns the new sprite position (X,Y).
func (ev *EventSpriteMove) Position() (int, int) {
	return ev.x, ev.y
}

// EventSpriteAccelerate is delivered ot the watchers when the sprite
// velocity changes.
type EventSpriteAccelerate struct {
	when time.Time
	s    *Sprite
	x, y float64
}

// When returns the time of the event.
func (ev *EventSpriteAccelerate) When() time.Time {
	return ev.when
}

// Sprite returns the sprite that accelerated.
func (ev *EventSpriteAccelerate) Sprite() *Sprite {
	return ev.s
}

// Position returns the new sprite velocity (X,Y).
func (ev *EventSpriteAccelerate) Velocity() (float64, float64) {
	return ev.x, ev.y
}
