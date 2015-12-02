.PHONY: server

GO=$(shell command -v go)
FSWATCH=$(shell command -v fswatch)
WATCH_GO=$(FSWATCH) -E --exclude=".*" --include="\.(tmpl|go)$$" ./

default: test build

test:
	$(GO) test ./...

build: install

install:
	$(GO) install -v ./...

server: install
	gocard server

autotest:
	$(WATCH_GO) | xargs -n1 -I{} $(MAKE) test

autobuild:
	$(WATCH_GO) | xargs -n1 -I{} $(MAKE) build

compile_daemon:
	CompileDaemon -exclude-dir=.git -command="gocard server" -color -pattern "(.+\\.go|.+\\.c|.+\\.tmpl)$$"

install_daemon:
	$(GO) get github.com/githubnemo/CompileDaemon
