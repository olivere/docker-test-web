PROJECT=github.com/olivere/docker-test-web
BUILDTIME=`date -u '+%Y-%m-%dT%H:%M:%SZ'`
BUILDTAG=`git rev-parse --short HEAD`
VERSION?=latest

default: build

.PHONY: build
build:
	go build -o docker-test-web $(PROJECT)

.PHONY: container
container:
#ifndef VERSION
#	$(error Please specify VERSION to tag the docker container)
#endif
	docker run --rm -v "$$PWD":/go/src/$(PROJECT) -w /go/src/$(PROJECT) -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 golang:1.7 go build -a -installsuffix cgo -ldflags "-w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILDTIME) -X main.BuildTag=$(BUILDTAG)" -o docker-test-web $(PROJECT)
	docker build -t olivere/docker-test-web:$(VERSION) .
