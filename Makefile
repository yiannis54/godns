run:
	go run src/cmd/*.go

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -v -o builds/godns.sh src/cmd/*.go
