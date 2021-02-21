#!/bin/bash

go build ./client
go build ./server
go build -o build/_examples/basic ./_examples/basic
go build -o build/_examples/echo_client ./_examples/echo/client
go build -o build/_examples/echo_server ./_examples/echo/server
