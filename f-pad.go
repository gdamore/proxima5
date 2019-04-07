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

func makePad(level *Level, props GameObjectProps) error {

	x := props.PropInt("x", 0)
	y := props.PropInt("y", 0)

	frame := props.PropString("frame", "PAD")
	sname := props.PropString("sprite", "Pad")
	sprite := GetSprite(sname)

	// sprite starts at 0, 0, so we can get its width
	// we use this to resize "TARMAC" frames.
	_, _, w, _ := sprite.Bounds()
	w++
	if nw := props.PropInt("width", w); nw != w {
		sprite.Resize(nw, 1)
	}
	sprite.SetLayer(LayerPad)
	sprite.SetFrame(frame)
	sprite.SetPosition(x, y)
	level.AddSprite(sprite)
	return nil
}

func init() {
	RegisterGameObjectMaker("pad", makePad)
}
