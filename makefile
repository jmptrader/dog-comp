all:
	ctags -R
	find . -name *.go |xargs gofmt -w
	go build -gcflags "-N -l"
clean:
	rm dog-comp
