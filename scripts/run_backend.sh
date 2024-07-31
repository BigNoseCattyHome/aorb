#!/bin/bash

# 启动Consul代理
echo "Starting Consul agent..."
nohup consul agent -dev >> ./logs/consul.log 2>&1 &

# 等待Consul启动
sleep 3

# 启动各个Go服务并输出实际监听的端口
function start_service() {
    service_name=$1
    service_cmd=$2
    log_file=$3

    echo "Starting $service_name..."
    nohup $service_cmd >> $log_file 2>&1 &
    sleep 1

    # pid=$(pgrep -f "$service_cmd")
    # if [ -n "$pid" ]; then
    #     port=$(lsof -Pan -p $pid -i | grep LISTEN | awk '{print $9}')
    #     echo "$service_name is listening on port $port"
    # else
    #     echo "$service_name failed to start"
    # fi
}

start_service "API Gateway" "go run ./backend/api-gateway/main.go" "./logs/api-gateway.log"
start_service "Auth Service" "go run ./backend/go-services/auth/main.go" "./logs/auth.log"
start_service "Comment Service" "go run ./backend/go-services/comment/main.go" "./logs/comment.log"
start_service "Poll Service" "go run ./backend/go-services/poll/main.go" "./logs/poll.log"
start_service "User Service" "go run ./backend/go-services/user/main.go" "./logs/user.log"
# start_service "Vote Service" "go run ./backend/go-services/vote/main.go" "./logs/vote.log"

echo "All services started."
