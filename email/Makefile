
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get micro.dev/v4/cmd/protoc-gen-micro
.PHONY: api
api:
	protoc --openapi_out=. --proto_path=. proto/email.proto
.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/email.proto
	
.PHONY: build
build:
	go build -o email *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t email:latest
