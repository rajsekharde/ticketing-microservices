# Define the default target when running 'make' with no arguments
.DEFAULT_GOAL := help

# Prevents name conflicts with real files or folders
.PHONY: proto run-gateway help


# Find all .proto files in the proto/ directory and its subdirectories
PROTO_FILES := $(shell find proto -name "*.proto")

## proto: Compile all .proto files across all services
proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO_FILES)
	@echo "Go code generated from Protobuf files"


## build-gateway: Build docker image for gateway service
build-gateway:
	docker build -f gateway/Dockerfile -t rajsekhar05/ticketing-ms-gateway:latest .

## build-user: Build docker image for user service
build-user:
	docker build -f user/Dockerfile -t rajsekhar05/ticketing-ms-user:latest .

## build-event: Build docker image for event service
build-event:
	docker build -f event/Dockerfile -t rajsekhar05/ticketing-ms-event:latest .

## build-booking: Build docker image for booking service
build-booking:
	docker build -f booking/Dockerfile -t rajsekhar05/ticketing-ms-booking:latest .

## build-payment: Build docker image for payment service
build-payment:
	docker build -f payment/Dockerfile -t rajsekhar05/ticketing-ms-payment:latest .

## build-notification: Build docker image for notification service
build-notification:
	docker build -f notification/Dockerfile -t rajsekhar05/ticketing-ms-notification:latest .


## build-all: Build docker images for all services
build-all: build-gateway build-user build-event build-booking build-payment build-notification

## docker-run: Run all pre-built service images using docker compose
docker-run:
	docker compose up -d

## docker-stop: Stop and remove all service containers
docker-stop:
	docker compose down --volumes --remove-orphans


## help: Show available commands
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'