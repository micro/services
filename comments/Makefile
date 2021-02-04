
GOPATH:=$(shell go env GOPATH)
MODIFY=Mgithub.com/micro/micro/proto/api/api.proto=github.com/micro/micro/v3/proto/api

.PHONY: proto
proto:
    
	protoc --openapi_out=. --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/comments.proto

.PHONY: docs
docs:
    
	protoc --openapi_out=. --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/comments.proto
	@redoc-cli bundle api-comments.json

.PHONY: build
build: proto

	go build -o comments-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t comments-service:latest
