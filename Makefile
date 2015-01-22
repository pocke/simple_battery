build:
	go-bindata -o icons.go -nomemcopy assets/
	go build
