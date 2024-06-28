#!/bin/bash

# 定义脚本目录和项目根目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(realpath "$SCRIPT_DIR/../")"

# 定义 Proto 文件目录和输出目录
PROTO_DIR="$PROJECT_ROOT/idl/"
OUTPUT_DART_DIR="$PROJECT_ROOT/frontend/lib/generated/protos/"
OUTPUT_GO_DIR="$PROJECT_ROOT/backend/rpc/"

# 定义 Dart 和 Go 插件路径
PROTOC_GEN_DART="$HOME/.pub-cache/bin/protoc-gen-dart"
PROTOC_GEN_GO="protoc-gen-go"
PROTOC_GEN_GO_GRPC="protoc-gen-go-grpc"

# 确保输出目录存在
mkdir -p "$OUTPUT_DART_DIR"
mkdir -p "$OUTPUT_GO_DIR"

# 检查 Proto 文件目录是否存在
if [ ! -d "$PROTO_DIR" ]; then
  echo "错误：Proto 目录不存在：$PROTO_DIR"
  exit 1
fi

# 检查 Dart 插件是否存在且可执行
if [ ! -x "$PROTOC_GEN_DART" ]; then
  echo "错误：protoc-gen-dart 不存在或不可执行：$PROTOC_GEN_DART"
  echo "请确保已安装 protoc_plugin 并设置正确的路径"
  exit 1
fi

# 检查 Go 插件是否存在
if ! command -v "$PROTOC_GEN_GO" &> /dev/null || ! command -v "$PROTOC_GEN_GO_GRPC" &> /dev/null; then
  echo "错误：Go 插件不存在。请确保已安装 protoc-gen-go 和 protoc-gen-go-grpc"
  exit 1
fi

echo "PROJECT_ROOT: "$PROJECT_ROOT
echo "OUTPUT_DO_DIR: "$OUTPUT_GO_DIR
echo "PROTO_DIR: "$PROTO_DIR

# 生成 Dart 文件
protoc --plugin=protoc-gen-dart="$PROTOC_GEN_DART" \
       --experimental_allow_proto3_optional \
       --dart_out=grpc:"$OUTPUT_DART_DIR" \
       -I"$PROTO_DIR" \
       "$PROTO_DIR"/*.proto

# 生成 Go 文件
protoc --go_out="$OUTPUT_GO_DIR" --go-grpc_out="$OUTPUT_GO_DIR" \
       -I"$PROTO_DIR" \
       "$PROTO_DIR"/*.proto

if [ $? -eq 0 ]; then
  echo "Dart 和 Go 文件已成功生成"
else
  echo "生成文件时发生错误"
  exit 1
fi