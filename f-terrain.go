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

func makeTerrain(level *Level, props GameObjectProps) error {

	frame := props.PropString("frame", "TERRAIN")
	sname := props.PropString("sprite", "")
	if sname == "" {
		return nil
	}
	sprite := GetSprite(sname)
	sprite.SetLayer(LayerTerrain)
	sprite.SetFrame(frame)
	level.AddSprite(sprite)
	return nil
}

func init() {
	RegisterGameObjectMaker("terrain", makeTerrain)
}
