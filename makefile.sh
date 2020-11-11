#!/bin/bash

make_build() {
    go build ./client
    go build ./server
    go build -o build/_examples/basic ./_examples/basic
    go build -o build/_examples/echo_client ./_examples/echo/client
    go build -o build/_examples/echo_server ./_examples/echo/server
}

make_api() {
    cd api && protoc --proto_path=. --go_out=,paths=source_relative:. *.proto
}

make_client() {
    go run ./_examples/echo/client
}

make_server() {
    go run ./_examples/echo/server
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
"server")
    make_server
    ;;
esac
