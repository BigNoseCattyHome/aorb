#!/bin/bash

# 定义要扫描的端口
#ports=(8300 37000 37001 37002 37003 37004 37005 37006 37007 37008 37009)
ports=(8300 37000 37001 37002 37003 37004 37005 37006)

echo "37000--网关"
echo "37001--Auth"
echo "37002--User"
echo "37003--Comment"
echo "37004--Vote"
echo "37005--Poll"
echo "37006--Message"

# 循环扫描每个端口
for port in "${ports[@]}"; do
    # 使用lsof命令查找绑定到该端口的进程
    process=$(lsof -i :$port)
    
    if [ -n "$process" ]; then
        echo "端口 $port 被以下进程绑定："
        echo "$process"
        read -p "是否终止该进程？(y/n): " answer
        
        if [ "$answer" == "y" ]; then
            # 获取进程ID并终止它
            pid=$(echo "$process" | awk 'NR==2 {print $2}')
            kill -9 $pid
            echo "进程 $pid 已被终止。"
        else
            echo "未终止进程。"
        fi
    else
        echo "端口 $port 未被绑定。"
    fi
done