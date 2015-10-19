# !/bin/sh

# Build the asset constructor first
go build mkassets.go spritedata.go leveldata.go properties.go || exit 1

# Then build the assets
./mkassets -type level l-*.yml || exit 1
./mkassets -type sprite s-*.yml || exit 1

# Now build the program
go build .
