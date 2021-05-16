#!/bin/bash

cd xcloud/apis && protoc --proto_path=. --go_out=,paths=source_relative:. *.proto
