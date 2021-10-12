#!/bin/sh
protoc -I ../protos --go_opt=module=github.com/mcfizh/grpc-demo --go_out=plugins=grpc:. value.proto
