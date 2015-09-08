all:
	ctags -R
	gofmt -w *.go
	go build -gcflags "-N -l"
