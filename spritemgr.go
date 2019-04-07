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

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type SpriteManager struct {
	width    int
	height   int
	maxlayer int
	minlayer int
	view     *views.ViewPort
	sprites  map[int]map[*Sprite]struct{}
}

func NewSpriteManager(width, height int) *SpriteManager {
	mgr := &SpriteManager{}
	mgr.sprites = make(map[int]map[*Sprite]struct{})
	mgr.width = width
	mgr.height = height
	return mgr
}

func (m *SpriteManager) forAllReverse(fn func(*Sprite)) {
	for layer := m.maxlayer; layer >= m.minlayer; layer-- {
		if sprites, ok := m.sprites[layer]; ok {
			for s := range sprites {
				fn(s)
			}
		}
	}
}

func (m *SpriteManager) forAll(fn func(*Sprite)) {
	for layer := m.minlayer; layer <= m.maxlayer; layer++ {
		if sprites, ok := m.sprites[layer]; ok {
			for s := range sprites {
				fn(s)
			}
		}
	}
}

func (m *SpriteManager) Update(now time.Time) {

	m.forAll(func(s *Sprite) {
		s.Update(now)
	})

	// Check for collisions, and report them
	// Note: because we report the collision *after* a move has occurred,
	// any corrective move (e.g. reverse position and move back) needs to
	// be done in the event handlers.

	// We go in reverse order, because higher layers get more precedence.
	// This way a sprite sitting on top of the terrain gets its handling
	// done first.

	sprites := make(map[*Sprite]struct{})
	m.forAllReverse(func(s1 *Sprite) {
		sprites[s1] = struct{}{}
		m.forAllReverse(func(s2 *Sprite) {
			if _, ok := sprites[s2]; !ok {
				if CollidersCollide(s1, s2) {
					HandleCollision(s1, s2)
				}
			}
		})
	})
}

func (m *SpriteManager) SetView(v *views.ViewPort) {
	m.view = v
	v.SetContentSize(m.width, m.height, true)
}

func (m *SpriteManager) Draw(v views.View) {

	v.Clear()

	m.forAll(func(s *Sprite) { s.Draw(v) })
}

func (m *SpriteManager) Size() (int, int) {
	return m.width, m.height
}

func (m *SpriteManager) HandleEvent(ev tcell.Event) bool {

	for layer := m.minlayer; layer <= m.maxlayer; layer++ {
		if sprites, ok := m.sprites[layer]; ok {
			for s := range sprites {
				if s.HandleEvent(ev) {
					return true
				}
			}
		}
	}
	return false
}

func (m *SpriteManager) AddSprite(s *Sprite) {

	layer := s.Layer()
	if m.sprites[layer] == nil {
		m.sprites[layer] = make(map[*Sprite]struct{})
	}
	m.sprites[layer][s] = struct{}{}
	if layer < m.minlayer {
		m.minlayer = layer
	}
	if layer > m.maxlayer {
		m.maxlayer = layer
	}
}

func (m *SpriteManager) RemoveSprite(s *Sprite) {
	delete(m.sprites[s.Layer()], s)
}

func (m *SpriteManager) Reset() {
	m.forAll(m.RemoveSprite)
}
