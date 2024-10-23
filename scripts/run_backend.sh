#!/bin/bash

# 启动Consul代理
echo "Starting Consul agent..."
nohup consul agent -dev >> ./logs/consul.log 2>&1 &

# 等待Consul启动
sleep 1

# 启动各个Go服务并输出实际监听的端口
function start_service() {
    service_name=$1
    service_dir=$2
    service_cmd=$3
    log_file=$4
    sleep 0.5

    echo "Starting $service_name..."
    pushd $service_dir > /dev/null
    nohup $service_cmd >> $log_file 2>&1 &
    popd > /dev/null
}

start_service "API Gateway" "./backend/api-gateway" "go run main.go" "./api-gateway.log"
start_service "Auth Service" "./backend/services/auth" "go run main.go" "./auth.log"
start_service "Comment Service" "./backend/services/comment" "go run main.go" "./comment.log"
start_service "Poll Service" "./backend/services/poll" "go run main.go" "./poll.log"
start_service "User Service" "./backend/services/user" "go run main.go" "./user.log"
start_service "Vote Service" "./backend/services/vote" "go run main.go" "./vote.log"
start_service "Message Service" "./backend/services/message" "go run main.go" "./message.log"

echo "All services started."
echo "Check logs for more details."
