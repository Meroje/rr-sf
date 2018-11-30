.PHONY: deps clean build

deps:
	cd aws-handler && dep ensure -v

clean: 
	rm -f ./bin/main
	
build:
	GOOS=linux GOARCH=amd64 go build -o bin/main ./aws-handler/main.go
