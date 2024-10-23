# 使用 golang:1.22 构建镜像，进行编译阶段
FROM golang:1.22 AS builder

# 设置 GOPROXY 环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录为 /app
WORKDIR /app

# 将 go.mod 和 go.sum 文件复制到镜像中
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 将当前目录下的所有文件复制到镜像中的 /app 目录
COPY . .

# 构建可执行文件，将其命名为 auth
RUN go build -o auth main.go

# 使用 ubuntu:22.04 作为运行时的基础镜像
FROM ubuntu:22.04

# 设置工作目录为 /app
WORKDIR /app

# 从构建阶段复制生成的二进制文件到运行时镜像
COPY --from=builder /app/auth .

# 暴露服务端口 8080
EXPOSE 8080

# 指定启动命令
CMD ["./auth"]