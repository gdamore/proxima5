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
)

type Collider interface {
	// Collides returns true if the Collider would collide with an
	// object at the given coordinates.
	Collides(x, y int) bool

	// Bounds is used to report the boundary for the Collider.  We
	// make use of this to optimize our checks, so we begin looking
	// at the more tightly constrained object first.  We also are
	// able to quickly eliminate two objects that do not overlap.
	// If an object never collides with anything, it sould specify
	// the value -1 for these coordinates. Nothing collides at any -1
	// coordinate, even with other objects also at -1 coordinates.
	Bounds() (x1, y1, x2, y2 int)

	// GetLayer returns the layer that the collider exists at.
	// This is mostly useful to class collisions as a group, e.g.
	// all hazards can be set to LayerHazard.
	Layer() int

	EventHandler
}

// CollidersCollide returns true if the two colliders collide.
func CollidersCollide(a, b Collider) bool {
	ax1, ay1, ax2, ay2 := a.Bounds()
	bx1, by1, bx2, by2 := b.Bounds()

	// Clip to the intersection of the two objects, and check
	// to make sure we still have overlap in each dimension.

	if bx1 > ax1 {
		ax1 = bx1
	}
	if bx2 < ax2 {
		ax2 = bx2
	}
	if ax1 > ax2 {
		return false
	}

	if by1 > ay1 {
		ay1 = by1
	}
	if by2 < ay2 {
		ay2 = by2
	}
	if ay1 > ay2 {
		return false
	}

	for y := ay1; y <= ay2; y++ {
		for x := ax1; x <= ax2; x++ {
			if a.Collides(x, y) && b.Collides(x, y) {
				return true
			}
		}
	}
	return false
}

// EventCollision is delivered to a Collider when a collision with another
// object is detected.  Both Colliders will receive an event, but it will
// not be the same event, since they see each other's collision.
type EventCollision struct {
	when     time.Time
	collider Collider
	target   Collider
}

// When returns the time when the collision occurred.
func (ev *EventCollision) When() time.Time {
	return ev.when
}

// Collider returns the object impacting the notified object.
func (ev *EventCollision) Collider() Collider {
	return ev.collider
}

// Target returns the object that is being hit.
func (ev *EventCollision) Target() Collider {
	return ev.target
}

// HandleCollision generates two EventCollisions, one for the each Collider,
// and deliveres them to the two Colliders.
func HandleCollision(c1, c2 Collider) {
	when := time.Now()
	ev1 := &EventCollision{when: when, collider: c2, target: c1}
	ev2 := &EventCollision{when: when, collider: c1, target: c2}

	c1.HandleEvent(ev1)
	c2.HandleEvent(ev2)
}
