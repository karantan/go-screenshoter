.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/screenshoter functions/screenshoter/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

deploy-fn: clean build
	sls deploy function --function screenshot --verbose
