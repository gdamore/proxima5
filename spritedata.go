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

type SpriteGlyph struct {
	Display    string // What is displayed, can be Unicode
	Foreground string // See colors.go for color names
	Background string
}

type SpriteFrameData struct {
	Names []string // names for frame
	Next  string   // Next Frame number, or "" for none
	Time  int      // time in msec to display frame
	Data  []string // lines of frame data (ASCII)
}

type SpriteData struct {
	Name    string
	Width   int
	Height  int
	OriginX int
	OriginY int
	Layer   int
	Glyphs  map[string]SpriteGlyph
	Palette map[string]string
	Frames  []SpriteFrameData
}
