SHELL=/bin/bash -o pipefail

FIRST_GOPATH:=$(firstword $(subst :, ,$(shell go env GOPATH)))
RR_BINARY:=$(FIRST_GOPATH)/bin/rr

bin/aws-handler: build

.PHONY: clean build

clean: 
	rm -f ./bin/aws-handler

build: vendor
	GOOS=linux GOARCH=amd64 go build -o bin/aws-handler ./bin/aws-handler.go

############
# Binaries #
############

vendor: tools.go composer.lock
	go mod vendor
	composer install

.PHONY: binaries
binaries: $(RR_BINARY)

$(RR_BINARY): vendor
	@go install -mod=vendor github.com/spiral/roadrunner/cmd/rr
