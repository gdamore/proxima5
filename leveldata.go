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

// LevelData describes the set up for the level.
//
// In order to keep this as flexible as possible, we use dictionaries of
// name value pairs.  The level has some initial properties that are
// global to the level itself, and then a series of objects to create.
// One of those objects must be a player object.
//
// Here are the properties we known about for the level (the Properties
// member of the LevelData):
//
//   title       Printable title for level
//   terrain     Name of terrain sprite.  The terrain sprite is expected
//               to have a starting frame called "TERRAIN".  This is
//               likely the only frame in the terrain sprite.
//   player      Name of player object (found in Objects array)
//   time        Timer in seconds for completion of the level.  Use
//               0 to create a non-expiring level.
//   gravity     (cells/sec*sec) This value is added to Y velocity.
//
// Each game object may have the following standard properties.  Different
// classes of objects may have others.  These are used for the Objects
// member.
//
//   label       Label for the instance.  This string value is set implictly
//               to the key of map, and it is not overridable.
//   class       Name of the object class (string).  This is used to obtain
//               an instance of the object, with MakeGameObject().
//               Note that the detailed list of class-specific bheaviors
//               and properties, as well as any requirements, are part
//               of each class.  For example, many classes require that
//               the sprites with which they are used meet certain
//               requirements.
//   sprite      Name of the sprite to use (string).  Each class provides
//               a default value here, but a different sprite can be
//               chosen with this.  (Not all classes support changing
//               the sprite.)
//   x           X position (int) (using the sprite's X origin)
//   y           Y position (int) (using the sprite's Y origin)
//   vx          X velocity (float64), cells/sec.  Negative values are
//               motion to the left.
//   vy          Y velocity (float64), cells/sec.  Negative values are
//               motion up.
//   width       Width of sprite in cells (int) if adjustable.
//   height      Height of sprite in cells (int) if adjustable.
//   colorX      Name of color (string) to substitute for given color
//               in the sprite's color map.
//
type LevelData struct {
	Name       string
	Properties GameObjectProps
	Objects    map[string]GameObjectProps
}
