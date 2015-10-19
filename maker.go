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
	"errors"
	"log"
	"sync"
)

// ObjectMaker is a function that creates game objects as described
// by the properties.  There may actually be multiple objects or composite
// objects taht result.  The ObjectMaker will register any created objects
// with the Level.
type GameObjectMaker func(l *Level, props GameObjectProps) error

var makerLk sync.Mutex
var objectMakers map[string]GameObjectMaker

func RegisterGameObjectMaker(name string, maker GameObjectMaker) {

	makerLk.Lock()
	if objectMakers == nil {
		objectMakers = make(map[string]GameObjectMaker)
	}
	objectMakers[name] = maker
	makerLk.Unlock()
}

func MakeGameObject(l *Level, name string, props GameObjectProps) error {
	makerLk.Lock()
	m, ok := objectMakers[name]
	makerLk.Unlock()
	if !ok || m == nil {
		log.Printf("No object maker for %s", name)
		return errors.New("Do not know how to make named object")
	}
	return m(l, props)
}
