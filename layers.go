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

const (
	LayerBackdrop  = 0 // Any backdrop (will we have?)
	LayerThrust    = 1 // Ship thrust (set low to hide)
	LayerTerrain   = 2 // Foreground terrain
	LayerHazard    = 3 // Hazards that the player can hit
	LayerShot      = 4 // Player's shots
	LayerPad       = 5 // Launch pad
	LayerPlayer    = 6 // The player's ship
	LayerExplosion = 7 // Player explosion
	LayerWipe      = 8 // Blastwave
	LayerDialog    = 9 // On top of everything
)
