all:
	go build


export:
	GOOS=linux GOARCH=arm go build && scp Oon oon:Oon/
#	GOOS=linux GOARCH=arm go test -c
