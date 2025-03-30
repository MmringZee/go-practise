#!/bin/bash

# 获取脚本所在目录为根路径
PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")

# 定义构建产物的输出目录为项目根目录下的_output文件夹
OUTPUT_DIR=${PROJ_ROOT_DIR}/_output

go build -o ${OUTPUT_DIR}/fg-apiserver -v cmd/fg-apiserver/main.go