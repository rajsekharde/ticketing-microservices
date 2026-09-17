# Define the default target when running 'make' with no arguments
.DEFAULT_GOAL := help

.PHONY: proto run-gateway help


# API Gateway server port
GATEWAY_PORT ?= ""

## run-gateway: Run the API Gateway
run-gateway:
	go run ./gateway -port ${GATEWAY_PORT}


# Find all .proto files in the proto/ directory and its subdirectories
PROTO_FILES := $(shell find proto -name "*.proto")

## proto: Compile all .proto files across all services
proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO_FILES)
	@echo "Go code generated from Protobuf files"


## help: Show available commands
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'