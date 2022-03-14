
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go install github.com/golang/protobuf/protoc-gen-go@latest
	go install github.com/micro/micro/v3/cmd/protoc-gen-micro@latest
	go install github.com/micro/micro/v3/cmd/protoc-gen-openapi@latest

.PHONY: api
api:
	protoc --openapi_out=. --proto_path=. proto/bitcoin.proto

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/bitcoin.proto
	
.PHONY: build
build:
	go build -o bitcoin *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t bitcoin:latest
