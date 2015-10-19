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
)

type TextBar struct {
	view        View
	leftStyle   tcell.Style
	rightStyle  tcell.Style
	centerStyle tcell.Style
	padStyle    tcell.Style
	left        []rune
	right       []rune
	center      []rune
}

func (t *TextBar) Draw() {
	v := t.view
	if v == nil {
		return
	}
	width, _ := v.Size()
	x := 0
	y := 0
	if width == 0 {
		return
	}

	v.Clear()

	for i := 0; i < width; i++ {
		v.SetContent(i, y, ' ', nil, t.padStyle)
	}

	// do left text
	for _, l := range t.left {
		v.SetContent(x, y, l, nil, t.leftStyle)
		x++
	}
	// advance for center if there is space
	if start := (width - len(t.center)) / 2; start > x {
		x = start
	}

	// do center text
	for _, l := range t.center {
		v.SetContent(x, y, l, nil, t.centerStyle)
		x++
	}

	// advance for right if there is space
	if start := width - len(t.right); start > x {
		x = start
	}

	// do right text
	for _, l := range t.right {
		v.SetContent(x, y, l, nil, t.rightStyle)
		x++
	}
}

func (t *TextBar) SetView(view View) {
	t.view = view
}

func (t *TextBar) HandleEvent(tcell.Event) bool {
	return false
}

func (t *TextBar) SetCenter(s string, style tcell.Style) {
	t.center = []rune(s)
	if style != tcell.StyleDefault {
		t.centerStyle = style
	}
}

func (t *TextBar) SetLeft(s string, style tcell.Style) {
	t.left = []rune(s)
	if style != tcell.StyleDefault {
		t.leftStyle = style
	}
}

func (t *TextBar) SetRight(s string, style tcell.Style) {
	t.right = []rune(s)
	if style != tcell.StyleDefault {
		t.rightStyle = style
	}
}

func (t *TextBar) SetStyle(style tcell.Style) {
	t.padStyle = style
}

func (t *TextBar) Resize() {
	// Nothing we can do.. move on.
}

func NewTextBar() *TextBar {
	t := &TextBar{
		left:   []rune{},
		right:  []rune{},
		center: []rune{},
	}
	return t
}
