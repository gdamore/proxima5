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

// EventHandler is anything that handles events.
type EventHandler interface {
	HandleEvent(tcell.Event) bool
}

type EventGravity struct {
	when  time.Time
	accel float64
}

func (ev *EventGravity) When() time.Time {
	return ev.when
}

func (ev *EventGravity) Accel() float64 {
	return ev.accel
}

type EventGame struct {
	when time.Time
}

func (ev *EventGame) When() time.Time {
	return ev.when
}

type EventPlayerDeath struct {
	EventGame
}

type EventLevelComplete struct {
	EventGame
}

type EventGameOver struct {
	EventGame
}

type EventLevelStart struct {
	EventGame
}

type EventTimesUp struct {
	EventGame
}
