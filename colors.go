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

var Colors16 = map[string]tcell.Color{
	"none":          tcell.ColorDefault,
	"black":         tcell.ColorBlack,
	"red":           tcell.ColorMaroon,
	"green":         tcell.ColorGreen,
	"yellow":        tcell.ColorOlive,
	"brown":         tcell.ColorOlive,
	"orange":        tcell.ColorOlive,
	"blue":          tcell.ColorNavy,
	"magenta":       tcell.ColorPurple,
	"cyan":          tcell.ColorTeal,
	"lightgrey":     tcell.ColorSilver,
	"grey":          tcell.ColorSilver,
	"darkgrey":      tcell.ColorSilver,
	"brightred":     tcell.ColorRed,
	"brightgreen":   tcell.ColorLime,
	"brightyellow":  tcell.ColorYellow,
	"brightblue":    tcell.ColorBlue,
	"brightmagenta": tcell.ColorFuchsia,
	"brightcyan":    tcell.ColorAqua,
	"white":         tcell.ColorWhite,
	"pink":          tcell.ColorFuchsia,
}

// On terminals with 256 colors, these will give a
// truer color representation.
var Colors256 = map[string]tcell.Color{
	"none":          tcell.ColorDefault,
	"black":         tcell.Color(232),
	"red":           tcell.Color(124),
	"green":         tcell.Color(34),
	"yellow":        tcell.Color(220),
	"orange":        tcell.Color(214),
	"brown":         tcell.Color(94),
	"blue":          tcell.Color(21),
	"magenta":       tcell.Color(163),
	"cyan":          tcell.Color(44),
	"lightgrey":     tcell.Color(251),
	"grey":          tcell.Color(245),
	"darkgrey":      tcell.Color(240),
	"brightred":     tcell.Color(196),
	"brightgreen":   tcell.Color(46),
	"brightyellow":  tcell.Color(226),
	"brightblue":    tcell.Color(27),
	"brightmagenta": tcell.Color(201),
	"brightcyan":    tcell.Color(51),
	"white":         tcell.Color(255),
	"pink":          tcell.Color(213),
}
