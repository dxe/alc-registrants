.PHONY: build deploy

build:
	go fmt ./...
	env GOOS=linux GOARCH=arm64 go build --ldflags="-s -w" -tags lambda.norpc -o bin/totalRegistrants/bootstrap src/totalRegistrants.go
	zip -j bin/totalRegistrants.zip bin/totalRegistrants/bootstrap

clean:
	rm -rf bin/*

deploy: clean build
	sls deploy --verbose
