#!/bin/bash

cd api && protoc --proto_path=. --go_out=,paths=source_relative:. *.proto
