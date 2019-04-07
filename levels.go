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
	"sync"
)

var levelsD map[string]*LevelData
var levelsL sync.Mutex

func GetLevel(name string) *Level {
	levelsL.Lock()
	defer levelsL.Unlock()
	if data, ok := levelsD[name]; ok {
		return NewLevelFromData(data)
	}
	return nil
}

func addLevel(data *LevelData) {
	levelsL.Lock()
	defer levelsL.Unlock()
	if levelsD == nil {
		levelsD = make(map[string]*LevelData)
	}
	levelsD[data.Name] = data
}

func RegisterLevelGobZ(b []byte) {
	r, e := gzip.NewReader(bytes.NewBuffer(b))
	if e != nil {
		panic("Malformed GZIP level data: " + e.Error())
	}
	dec := gob.NewDecoder(r)
	data := &LevelData{}
	if e := dec.Decode(data); e != nil {
		panic("Malformed GOB level data: " + e.Error())
	}
	addLevel(data)
}
