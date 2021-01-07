GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/locations.proto
	
.PHONY: build
build:
	go build -o locations *.go

.PHONY: test
test:
	go test -v ./... -cover