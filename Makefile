buid:
	go build

autobuild:
	CompileDaemon -exclude-dir=.git -command="./gocard server" -color -pattern "(.+\\.go|.+\\.c|.+\\.tmpl)$$"

install_daemon:
	go get github.com/githubnemo/CompileDaemon

server: build
	@./gocard
