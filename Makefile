
.PHONY: build
build:
	@go build .
	@go build -o build/_examples/basic ./_examples/basic
	@go build -o build/_examples/echo_server ./_examples/echo/server

.PHONY: api
api:
	@cd api && protoc --proto_path=. --go_out=,paths=source_relative:. *.proto