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
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gdamore/proxima5"
	"gopkg.in/yaml.v2"
)

// This program builds sprite data as go declarations, loading the data
// from YAML.

func die(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	fmt.Fprintln(os.Stderr, "")
	os.Exit(1)
}

func validateSprite(iname string, data *proxima5.SpriteData) {
	if data.Name == "" {
		die("%s: Missing name", iname)
	}

	for fri, frame := range data.Frames {
		if len(frame.Data) != data.Height {
			die("Frame %d bad lines (%d != %d)",
				fri, len(frame.Data), data.Height)
		}
		for lno, line := range frame.Data {
			if len(line) != data.Width {
				die("%s: Frame %d, line %d wrong len "+
					"(%d != %d)", iname, fri, lno,
					len(line), data.Width)
			}
			for i := range line {
				if line[i] == ' ' {
					continue
				}
				ss := string(line[i])
				if _, ok := data.Glyphs[ss]; !ok {
					die("%s: Frame %d: line %d: "+
						"unknown glyph",
						iname, fri, i)
				}
			}
		}
	}
}

func writeSpriteGo(w io.Writer, data *proxima5.SpriteData, pkg string) {
	buf := new(bytes.Buffer)
	zbuf, e := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if e != nil {
		die("Cannot gzip gob: %v", e)
	}
	enc := gob.NewEncoder(zbuf)
	if e := enc.Encode(data); e != nil {
		die("Cannot encode gob: %v", e)
	}
	zbuf.Close()

	fmt.Fprintf(w, "// Sprite data file\n")
	fmt.Fprintf(w, "// Generated automatically, do not edit!\n")
	fmt.Fprintf(w, "package %s\n\n", pkg)

	fmt.Fprintf(w, "func init() {\n")
	fmt.Fprintf(w, "\tRegisterSpriteGobZ([]byte{\n")
	x := buf.Bytes()
	for i, b := range x {
		if i%16 == 0 {
			fmt.Fprint(w, "\t\t")
		} else {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, "0x%02x,", b)
		if i%16 == 15 {
			fmt.Fprint(w, "\n")
		}
	}
	if len(x)%16 != 0 {
		fmt.Fprint(w, "\n")
	}
	fmt.Fprint(w, "\t})\n}\n")
}

func validateLevel(iname string, level *proxima5.LevelData) {
	if level.Name == "" {
		die("Missing level label")
	}
	if level.Properties.PropInt("width", 0) < 1 ||
		level.Properties.PropInt("height", 0) < 1 {
		die("Impossible level dimensions")
	}
}

func writeLevelGo(w io.Writer, level *proxima5.LevelData, pkg string) {

	buf := new(bytes.Buffer)
	zbuf, e := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if e != nil {
		die("Cannot gzip gob: %v", e)
	}
	enc := gob.NewEncoder(zbuf)
	if e := enc.Encode(level); e != nil {
		die("Cannot encode gob: %v", e)
	}
	zbuf.Close()

	fmt.Fprintf(w, "// Level data file\n")
	fmt.Fprintf(w, "// Generated automatically, do not edit!\n")
	fmt.Fprintf(w, "package %s\n\n", pkg)

	fmt.Fprintf(w, "func init() {\n")
	fmt.Fprintf(w, "\tRegisterLevelGobZ([]byte{\n")
	x := buf.Bytes()
	for i, b := range x {
		if i%16 == 0 {
			fmt.Fprint(w, "\t\t")
		} else {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, "0x%02x,", b)
		if i%16 == 15 {
			fmt.Fprint(w, "\n")
		}
	}
	if len(x)%16 != 0 {
		fmt.Fprint(w, "\n")
	}
	fmt.Fprint(w, "\t})\n}\n")
}

func main() {
	var (
		inf                       *os.File
		e                         error
		asset, oname, format, pkg string
	)

	flag.StringVar(&asset, "type", "sprite", "asset type [sprite|level]")
	flag.StringVar(&pkg, "pkg", "proxima5", "package name")
	flag.StringVar(&oname, "out", "", "output file name")
	flag.StringVar(&format, "format", "go", "output format [go|gob]")

	flag.Parse()

	for _, arg := range flag.Args() {

		data := &proxima5.SpriteData{}
		level := &proxima5.LevelData{}

		iname := arg
		if inf, e = os.Open(iname); e != nil {
			die("%s: open: %v", iname, e)
		}

		nname := oname
		if nname == "" {
			nname = iname + ".go"
			if strings.HasSuffix(iname, ".yml") {
				nname = iname[:len(iname)-4] + ".go"
			} else if strings.HasSuffix(iname, ".yaml") {
				nname = iname[:len(iname)-5] + ".go"
			}
		}
		all, e := ioutil.ReadAll(inf)
		if e != nil {
			die("%s; Cannot load: %v", iname, e)
		}

		buf := new(bytes.Buffer)

		switch asset {
		case "sprite":

			if err := yaml.Unmarshal(all, &data); err != nil {
				die("%s: YAML error: %v", iname, err)
			}
			validateSprite(iname, data)
			writeSpriteGo(buf, data, pkg)

		case "level":
			if err := yaml.Unmarshal(all, &level); err != nil {
				die("%s: YAML error: %v", iname, err)
			}
			validateLevel(iname, level)
			writeLevelGo(buf, level, pkg)

		default:
			panic("unknown asset class")
		}

		allb := buf.Bytes()
		if nname != "" && nname != "-" {
			e := ioutil.WriteFile(nname, allb, 0644)
			if e != nil {
				die("Cannot create output %s: %v", nname, e)
			}
		} else {
			os.Stdout.Write(allb)
		}
	}
}
