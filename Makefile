.PHONY: deps clean build

deps:
	go get -u

clean: 
	rm -f ./bin/aws-handler
	
build:
	GOOS=linux GOARCH=amd64 go build -o bin/aws-handler ./bin/aws-handler.go
