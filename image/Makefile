
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
	go get github.com/micro/micro/v3/cmd/protoc-gen-openapi
.PHONY: proto
proto:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=. proto/image.proto

.PHONE: api
api:
	protoc --openapi_out=. --proto_path=. proto/image.proto
	
.PHONY: build
build:
	go build -o image *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t image:latest
