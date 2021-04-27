#!/bin/bash

if [ $# -eq 0 ]; then
    echo "must have command"
    exit 1
fi

cmd=./scripts/makefiles/$1.sh
if [ -f $cmd ]; then
    sh $cmd ${@:2}
else
    echo "unsupport command $1"
    exit 2
fi
