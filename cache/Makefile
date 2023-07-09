
GOPATH:=$(shell go env GOPATH)

.PHONY: api
api:
	protoc --openapi_out=. --proto_path=. proto/cache.proto

.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get micro.dev/v4/cmd/protoc-gen-micro
	go get micro.dev/v4/cmd/protoc-gen-openapi

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/cache.proto
	
.PHONY: build
build:
	go build -o cache *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t cache:latest
