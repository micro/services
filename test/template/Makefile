
GOPATH:=$(shell go env GOPATH)
.PHONY: proto
proto:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=:. proto/idiomatic.proto

.PHONY: docs
docs:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=:. proto/idiomatic.proto
	@redoc-cli bundle api-protobuf.json

.PHONY: build
build:
	go build -o idiomatic *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t idiomatic:latest
