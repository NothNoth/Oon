all:
	go build

test:
	go test

export:
	GOOS=linux GOARCH=arm go build
#	GOOS=linux GOARCH=arm go test -c
