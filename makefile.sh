#!/bin/bash

make_build() {
    go build .
    go build -o build/_examples/basic ./_examples/basic
    go build -o build/_examples/echo_server ./_examples/echo/server
}

make_api() {
    cd api && protoc --proto_path=. --go_out=,paths=source_relative:. *.proto
}

make_client() {
    go run ./_examples/echo/client
}

case $1 in
"api")
    make_api
    ;;
"build")
    make_build
    ;;
"client")
    make_client
    ;;
esac
