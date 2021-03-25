
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
.PHONY: proto
proto:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=. proto/chats.proto

docs:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=. proto/chats.proto
	@redoc-cli bundle api-chats.json

.PHONY: build
build:
	go build -o chats *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t chats:latest
