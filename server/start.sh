#!/bin/bash
cd $(dirname $0)

server_name=passwordserver
configPath=./deploy/config/config.yaml

# 运行命令
./${server_name} start --config ${configPath}