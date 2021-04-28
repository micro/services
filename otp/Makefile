
GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v3/cmd/protoc-gen-micro
.PHONY: proto
proto:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=. proto/otp.proto

docs:
	protoc --openapi_out=. --proto_path=. --micro_out=. --go_out=. proto/otp.proto
	@redoc-cli bundle api-otp.json
	
.PHONY: build
build:
	go build -o otp *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t otp:latest
