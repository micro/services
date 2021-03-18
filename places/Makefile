GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro

.PHONY: proto
proto:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=:. proto/places.proto

.PHONY: docs
docs:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=:. proto/places.proto
	@redoc-cli bundle api-places.json

.PHONY: build
build:
	go build -o places *.go

.PHONY: test
test:
	go test -v ./... -cover
