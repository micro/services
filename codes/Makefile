
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/codes.proto
	
.PHONY: build
build:
	go build -o codes *.go

.PHONY: docs
docs:
        protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=:. proto/codes.proto
        @redoc-cli bundle api-codes.json

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t codes:latest
