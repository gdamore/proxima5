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
	"time"
)

// Alarm represents an alarm event.
type Alarm struct {
	expire time.Time
	h      EventHandler
}

// AlarmClock manages Alarms.  Each Alarm is set on a clock, and a
// single AlarmClock can have unlimited alarms registered with it.
// However, the actual alarms only ring (run the handlers) when the
// AlarmClock's Tick method is executed; no implicit asynchronous
// handling is performed.
type AlarmClock map[*Alarm]struct{}

// NewAlarmClock creates an AlarmClock.
func NewAlarmClock() AlarmClock {
	return AlarmClock{}
}

// Schedule schedules an alarm to fire after the duration has expired.
func (c AlarmClock) Schedule(d time.Duration, h EventHandler) *Alarm {
	a := &Alarm{expire: time.Now().Add(d), h: h}
	c[a] = struct{}{}
	return a
}

// Cancel cancels a previously scheduled alarm.
func (c AlarmClock) Cancel(a *Alarm) {
	delete(c, a)
}

// Tick runs through the list of alarms on this clock,, and runs the
// event handler for any that have expired.  Any alarm that is fired
// is descheduled.  Note that handlers all run synchornously within
// the caller's context.
func (c AlarmClock) Tick(now time.Time) {
	for a := range c {
		if a.expire.Before(now) {
			delete(c, a)
			ev := &EventAlarm{when: now}
			a.h.HandleEvent(ev)
		}
	}
}

// Reset resets the clock, clearing all alarms.
func (c AlarmClock) Reset() {
	for a := range c {
		delete(c, a)
	}
}

// EventAlarm is dispatched to the handler for an Alarm.
type EventAlarm struct {
	when time.Time
}

// When returns the time when the alarm fired.
func (ev *EventAlarm) When() time.Time {
	return ev.when
}
