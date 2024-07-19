.PHONY: build clean deploy remove

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bootstrap consumer/main.go

clean:
	rm -rf ./bootstrap

deploy: clean build
	sls deploy --verbose

remove:
	sls remove --verbose