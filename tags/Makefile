
GOPATH:=$(shell go env GOPATH)
MODIFY=Mgithub.com/micro/micro/proto/api/api.proto=github.com/micro/micro/v3/proto/api

.PHONY: proto
proto:
    
	protoc --openapi_out=. --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/tags.proto

.PHONY: docs
docs:
    
	protoc --openapi_out=. --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/tags.proto
	@redoc-cli bundle api-tags.json

.PHONY: build
build: proto

	go build -o tag-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t tag-service:latest
