// Copyright 2016 Garrett D'Amore
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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

const Version = "0.2"

// Proxima (Escape from Promixa 5) is a text-oriented game, utilizing
// Unicode characters and text terminal display capabilties, demonstrating
// the capabilities of tcell.   It needs at least an 80x24 column terminal
// to run well, and the highest baud rade you can arrange.

func main() {

	var logfile string

	flag.StringVar(&logfile, "log", logfile, "Log file for debugging log")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	if logfile != "" {
		if f, e := os.Create(logfile); e == nil {
			log.SetOutput(f)
		}
	} else {
		log.SetOutput(ioutil.Discard)
	}
	game := &Game{}
	if err := game.Init(); err != nil {
		fmt.Printf("Failed to initialize game: %v\n", err)
		os.Exit(1)
	}
	if err := game.Run(); err != nil {
		fmt.Printf("Failed to run game: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Escape from Proxima 5, %s\n", Version)
	fmt.Printf("Copyright 2016 Garrett D'Amore\n")
	fmt.Printf("Licensed under the Apache 2 license.\n")
	fmt.Printf("See http://www.apache.org/licenses/LICENSE-2.0\n")
	fmt.Printf("Thanks for playing!\n")
}
