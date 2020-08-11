.PHONY: build build-darwin build-amd64-linux

build: build-darwin build-amd64-linux

build-darwin:
	GOOS=darwin go build -o dist/httpeek-darwin httpeek.go

build-amd64-linux:
	GOOS=linux GOARCH=amd64 go build -o dist/httpeek-amd64-linux httpeek.go
