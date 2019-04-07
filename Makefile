.PHONY: assets run clean

DESTDIR ?=
PREFIX ?= /usr

cmd/proxima5/proxima5: assets
	@(cd cmd/proxima5 && go build)
	@echo 'Execute cmd/proxima5/proxima5 to launch the game'

assets: cmd/mkassets/mkassets
	@cmd/mkassets/mkassets -type level levels/l-*.yml
	@mv levels/*.go .
	@cmd/mkassets/mkassets -type sprite sprites/s-*.yml
	@mv sprites/*.go .

cmd/mkassets/mkassets:
	@(cd cmd/mkassets && go build)

install:
	@install -Dm755 cmd/proxima5/proxima5 $(DESTDIR)$(PREFIX)/bin/proxima5
	@install -Dm644 LICENSE $(DESTDIR)$(PREFIX)/share/licenses/proxima5/LICENSE

clean:
	@go clean
	@(cd cmd/proxima5 && go clean)
	@(cd cmd/mkassets && go clean)
	@rm -f l-*.go s-*.go
