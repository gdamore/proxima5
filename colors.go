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
	"red":           tcell.ColorRed,
	"green":         tcell.ColorGreen,
	"yellow":        tcell.ColorYellow,
	"brown":         tcell.ColorYellow,
	"orange":        tcell.ColorYellow,
	"blue":          tcell.ColorBlue,
	"magenta":       tcell.ColorMagenta,
	"cyan":          tcell.ColorCyan,
	"lightgrey":     tcell.ColorWhite,
	"grey":          tcell.ColorGrey,
	"darkgrey":      tcell.ColorGrey,
	"brightred":     tcell.ColorBrightRed,
	"brightgreen":   tcell.ColorBrightGreen,
	"brightyellow":  tcell.ColorYellow,
	"brightblue":    tcell.ColorBrightBlue,
	"brightmagenta": tcell.ColorBrightMagenta,
	"brightcyan":    tcell.ColorBrightCyan,
	"white":         tcell.ColorBrightWhite,
	"pink":          tcell.ColorBrightMagenta,
}

// On terminals with 256 colors, these will give a
// truer color representation.
var Colors256 = map[string]tcell.Color{
	"none":          tcell.ColorDefault,
	"black":         tcell.Color(232 + 1),
	"red":           tcell.Color(124 + 1),
	"green":         tcell.Color(34 + 1),
	"yellow":        tcell.Color(220 + 1),
	"orange":        tcell.Color(214 + 1),
	"brown":         tcell.Color(94 + 1),
	"blue":          tcell.Color(21 + 1),
	"magenta":       tcell.Color(163 + 1),
	"cyan":          tcell.Color(44 + 1),
	"lightgrey":     tcell.Color(251 + 1),
	"grey":          tcell.Color(245 + 1),
	"darkgrey":      tcell.Color(240 + 1),
	"brightred":     tcell.Color(196 + 1),
	"brightgreen":   tcell.Color(46 + 1),
	"brightyellow":  tcell.Color(226 + 1),
	"brightblue":    tcell.Color(27 + 1),
	"brightmagenta": tcell.Color(201 + 1),
	"brightcyan":    tcell.Color(51 + 1),
	"white":         tcell.Color(255 + 1),
	"pink":          tcell.Color(213 + 1),
}
