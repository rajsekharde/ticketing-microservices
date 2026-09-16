GATEWAY_PORT ?= ""

run-gateway:
	go run ./gateway -port ${GATEWAY_PORT}