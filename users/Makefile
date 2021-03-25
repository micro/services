
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
.PHONY: proto
proto:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=:. proto/users.proto

.PHONY: docs
docs:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=:. proto/users.proto
	@redoc-cli bundle api-users.json

.PHONY: build
build:
	go build -o users *.go

.PHONY: test
test:
	go test -v -p 1 ./...

.PHONY: docker
docker:
	docker build . -t users:latest
