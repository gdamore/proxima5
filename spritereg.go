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
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"strconv"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

var spriteData map[string]*SpriteData
var spriteDataL sync.Mutex

func addSprite(data *SpriteData) {
	spriteDataL.Lock()
	if spriteData == nil {
		spriteData = make(map[string]*SpriteData)
	}
	spriteData[data.Name] = data
	spriteDataL.Unlock()
}

func RegisterSpriteGobZ(b []byte) {
	r, e := gzip.NewReader(bytes.NewBuffer(b))
	if e != nil {
		panic("Malformed GZIP sprite data: " + e.Error())
	}
	dec := gob.NewDecoder(r)
	data := &SpriteData{}
	if e := dec.Decode(data); e != nil {
		panic("Malformed GOB sprite data: " + e.Error())
	}
	addSprite(data)
}

func GetSprite(name string) *Sprite {

	spriteDataL.Lock()
	data, ok := spriteData[name]
	spriteDataL.Unlock()
	if !ok {
		return nil
	}

	s := NewSprite(data.Width, data.Height)

	glyphs := make(map[byte]rune)
	styles := make(map[byte]tcell.Style)
	count := data.Width * data.Height

	for k, g := range data.Glyphs {
		st := tcell.StyleDefault
		fn := g.Foreground
		if cn, ok := data.Palette[fn]; ok {
			fn = cn
		}
		bn := g.Background
		if cn, ok := data.Palette[bn]; ok {
			bn = cn
		}

		fg := tcell.GetColor(fn)
		st = st.Foreground(fg)

		bg := tcell.GetColor(bn)
		st = st.Background(bg)

		glyphs[k[0]] = []rune(g.Display)[0]
		styles[k[0]] = st
	}

	for i, fr := range data.Frames {
		f := &spriteFrame{nextFrame: fr.Next}
		f.timer = time.Millisecond * time.Duration(fr.Time)
		f.runes = make([]rune, count)
		f.styles = make([]tcell.Style, count)
		for y, line := range fr.Data {
			for x := range line {
				c := line[x]
				i := y*data.Width + x
				if c == ' ' {
					continue
				}
				f.runes[i] = glyphs[c]
				f.styles[i] = styles[c]
			}
		}
		for _, name := range fr.Names {
			s.frames[name] = f
		}
		// Also insert a default based on index position
		name := strconv.Itoa(i)
		if _, ok := s.frames[name]; ok {
			s.frames[name] = f
		}
		s.frames[name] = f
	}
	s.SetOrigin(data.OriginX, data.OriginY)
	s.SetLayer(data.Layer)
	return s
}
