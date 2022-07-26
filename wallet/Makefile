
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro

.PHONY: api
api:
	protoc --openapi_out=. --proto_path=. proto/wallet.proto

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/wallet.proto

.PHONY: build
build:
	go build -o wallet *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t wallet:latest
