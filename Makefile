# Define the default target when running 'make' with no arguments
.DEFAULT_GOAL := help

# Prevents name conflicts with real files or folders
.PHONY: proto run-gateway help

## run-gateway: Start the API Gateway
run-gateway:
	go run ./gateway

## run-user: Start the User service
run-user:
	go run ./user

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