#!/bin/bash

dn=$(dirname $0)
cd $dn

echo "Checking for Go"
go version > /dev/null

if [ $? -ne 0 ]; then
	echo "Go is required to build this project, and the go"
	echo "executable must be on your PATH."
	echo 
	echo "Get Go from http://golang.org"
	exit 1
fi

echo "Getting dependencies"
# Make sure our dependencies are present
go get github.com/gdamore/tcell
go get gopkg.in/yaml.v2

echo "Building asset creator"
# Build the asset constructor first
go build mkassets.go spritedata.go leveldata.go properties.go || exit 1

echo "Generating assets"
./mkassets -type level l-*.yml || exit 1
./mkassets -type sprite s-*.yml || exit 1

echo "Building game"
# Now build the program
go build .

echo "Finished."
echo "Execute ${dn}/proxima5 to play.  Have fun!"
