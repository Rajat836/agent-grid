#!/usr/bin/env bash

# Prerequisite:
# brew install protobuf && protoc --version
# gsed or sed (brew install gnu-sed) is required for mac
# run the script from the root folder: ./scripts/gen_proto.sh
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# export PATH="$PATH:$(go env GOPATH)/bin"


sed="sed"
if [[ $OSTYPE == 'darwin'* ]]; then
  sed="gsed"
fi

#---------------------------Build start

PWD=$(pwd)
GO_OUT_FOLDER_NAME="$PWD/in@types"
NODE_OUT_FOLDER_NAME="$PWD/panel/@types"
PROTO_FOLDER="$PWD/@proto"

# if pwd ends with scripts, show error
if [[ $PWD == *"/scripts" ]]; then
  echo "Error: Please run the script from the root folder, not from the scripts folder."
  exit 1
fi

protoc \
  -I="proto" \
  --go_out="app/internal/types" \
  --go_opt=paths="source_relative" \
  $(find "proto" -name "*.proto")

# replace omitempty in go protos
find "app/internal/types" -type f -name "*.go" -exec sed -i "" -e "s/,omitempty//g" {} +

cd scripts && bun install && cd ..

protoc \
  -I="proto" \
  --ts_proto_out="panel/src/types" \
  --plugin="scripts/node_modules/.bin/protoc-gen-ts_proto" \
  --ts_proto_opt=snakeToCamel=false,forceLong=string,esModuleInterop=true,onlyTypes=true,useDate=false \
  $(find "proto" -name "*.proto")

echo "Generated files successfully"

# linting
cd panel && bun lint && cd ..
cd app && make lint && cd ..
