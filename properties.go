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
	"log"
	"strconv"
)

// GameObjectProps is a dictionary of properties for game objects.
// We create our own type because we have convenience methods that
// provide both type conversion and default values.
type GameObjectProps map[string]string

const (
	GameObjectPropX      = "x"       // int
	GameObjectPropY      = "y"       // int
	GameObjectPropWidth  = "width"   // int
	GameObjectPropHeight = "height"  // int
	GameObjectPropVX     = "vx"      // float64
	GameObjectPropVY     = "vy"      // float64
	GameObjectFrame      = "frame"   // start frame (string)
	GameObjectSprite     = "sprite"  // sprite name (string)
	GameObjectColor1     = "color1"  // color name (string)
	GameObjectColor2     = "color2"  // color name (string)
	GameObjectColor3     = "color3"  // color name (string)
	GameObjectColor4     = "color4"  // color name (string)
	GameObjectColor5     = "color5"  // color name (string)
	GameObjectColor6     = "color6"  // color name (string)
	GameObjectColor7     = "color7"  // color name (string)
	GameObjectColor8     = "color8"  // color name (string)
	GameObjectColor9     = "color9"  // color name (string)
	GameObjectColor10    = "color10" // color name (string)
	GameObjectColor11    = "color11" // color name (string)
	GameObjectColor12    = "color12" // color name (string)
	GameObjectColor13    = "color13" // color name (string)
	GameObjectColor14    = "color14" // color name (string)
	GameObjectColor15    = "color15" // color name (string)
	GameObjectColor16    = "color16" // color name (string)
)

func (p GameObjectProps) PropFloat64(name string, def float64) float64 {
	if s, ok := p[name]; !ok {
		return def
	} else if f, e := strconv.ParseFloat(s, 64); e != nil {
		log.Printf("Parse float64, prop %s: %v", name, e)
		return def
	} else {
		return f
	}
}

func (p GameObjectProps) PropInt(name string, def int) int {
	if s, ok := p[name]; !ok {
		return def
	} else if i, e := strconv.Atoi(s); e != nil {
		log.Printf("Parse int, prop %s: %v", name, e)
		return def
	} else {
		return i
	}
}

func (p GameObjectProps) PropBool(name string, def bool) bool {
	if s, ok := p[name]; !ok {
		return def
	} else if b, e := strconv.ParseBool(s); e != nil {
		log.Printf("Parse bool, prop %s: %v", name, e)
		return def
	} else {
		return b
	}
}

func (p GameObjectProps) PropString(name string, def string) string {
	if s, ok := p[name]; !ok {
		return def
	} else {
		return s
	}
}

func (p GameObjectProps) PropSetInt(name string, val int) {
	p[name] = strconv.Itoa(val)
}

func (p GameObjectProps) PropSetFloat64(name string, val float64) {
	p[name] = strconv.FormatFloat(val, 'f', 2, 64)
}

func (p GameObjectProps) PropSetBool(name string, val bool) {
	p[name] = strconv.FormatBool(val)
}

func (p GameObjectProps) PropSetString(name string, val string) {
	p[name] = val
}
