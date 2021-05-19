
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --openapi_out=. --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/helloworld.proto

.PHONY: api
api:
	protoc --openapi_out=. --proto_path=${GOPATH}/src:. proto/helloworld.proto

.PHONY: build
build: proto

	go build -o helloworld-srv *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t helloworld-srv:latest
