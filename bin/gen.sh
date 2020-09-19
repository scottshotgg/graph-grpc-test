#!/bin/bash

cd proto

# Look at: https://github.com/danielvladco/go-proto-gql


protoc \
  -I . \
  -I /usr/local/include \
  --go_out=plugins=grpc:../pkg/grapherino \
  *.proto