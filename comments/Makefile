
GOPATH:=$(shell go env GOPATH)

.PHONY: api
api:
	protoc --openapi_out=. --proto_path=. proto/comments.proto

.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
	go get github.com/micro/micro/v3/cmd/protoc-gen-openapi

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/comments.proto
	
.PHONY: build
build:
	go build -o comments *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t comments:latest
