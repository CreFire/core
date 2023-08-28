#!/usr/bin/env bash

 ./bin/protoc.exe --go_out=./pb/. --plugin=protoc-gen-go=./bin/protoc-gen-go.exe \
 --go-grpc_out=./pb/. --plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc.exe \
 ./pb/proto/*.proto