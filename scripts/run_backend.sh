#!/bin/bash

# 启动Consul代理
echo "Starting Consul agent..."
nohup consul agent -dev >> ./logs/consul.log 2>&1 &

# 等待Consul启动
sleep 1

# 启动各个Go服务并输出实际监听的端口
function start_service() {
    service_name=$1
    service_cmd=$2
    log_file=$3
    sleep 0.5

    echo "Starting $service_name..."
    nohup $service_cmd >> $log_file 2>&1 &
}

start_service "API Gateway" "go run ./backend/api-gateway/main.go" "./logs/api-gateway.log"
start_service "Auth Service" "go run ./backend/services/auth/main.go" "./logs/auth.log"
start_service "Comment Service" "go run ./backend/services/comment/main.go" "./logs/comment.log"
start_service "Poll Service" "go run ./backend/services/poll/main.go" "./logs/poll.log"
start_service "User Service" "go run ./backend/services/user/main.go" "./logs/user.log"
start_service "Vote Service" "go run ./backend/services/vote/main.go" "./logs/vote.log"
start_service "Message Service" "go run ./backend/services/message/main.go" "./logs/message.log"

echo "All services started."
echo "Check logs for more details."
