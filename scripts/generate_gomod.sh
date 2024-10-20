#! /bin/bash
# 在./backend/services/下的每一个子目录中生成go.mod文件
cd ../backend/services
for service in *; do
    if [ -d "$service" ]; then
        cd $service
        go mod init github.com/BigNoseCattyHome/aorb/backend/services/$service
        go mod tidy
        cd ..
    fi
done