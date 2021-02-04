
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/files.proto

.PHONY: docs
docs:    
	protoc --openapi_out=. --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/files.proto
	@redoc-cli bundle api-files.json

.PHONY: build
build:
	go build -o files *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t files:latest
