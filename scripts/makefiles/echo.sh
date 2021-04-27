#!/bin/bash

if [ $# -eq 0 ]; then
    echo "must spefic c/s"
elif [ $1 = "c" ]; then
    go run ./_examples/echo/client
elif [ $1 = "s" ]; then
    go run ./_examples/echo/server
else
    echo "must spefic c/s"
fi
