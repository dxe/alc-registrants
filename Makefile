.PHONY: build deploy

build:
	go fmt ./...
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/totalRegistrants src/totalRegistrants.go

clean:
	rm -rf bin/*

deploy: clean build
	sls deploy --verbose
