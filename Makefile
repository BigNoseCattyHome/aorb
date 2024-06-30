.PHONY: all clean test build deploy run_backend run_frontend proto

# 默认执行的命令
all: build proto

PROTO_PATH=./idl
GOOGLE_PROTO_PATH=/usr/local/include/google/protobuf
OUTPUT_DART_PATH=./frontend/lib/generated/
GO_OUT_PATH=./backend/rpc
PROTOC_GEN_DART=$(shell which protoc-gen-dart)

proto:
	@echo "Creating golang and dart grpc files..."
	@for file in $(PROTO_PATH)/*.proto; do \
        if [ -f "$$file" ]; then \
            prefix=$$(basename "$$file" .proto); \
            mkdir -p $(GO_OUT_PATH)/"$${prefix}"; \
            mkdir -p $(OUTPUT_DART_PATH); \
            echo "Created directory for $$prefix"; \
            protoc -I$(PROTO_PATH) \
                --go_out=$(GO_OUT_PATH)/$$prefix --go_opt=paths=source_relative \
                --go-grpc_out=$(GO_OUT_PATH)/$$prefix --go-grpc_opt=paths=source_relative \
                --dart_out=grpc:$(OUTPUT_DART_PATH) \
                --plugin=protoc-gen-dart=$(PROTOC_GEN_DART) \
                $$file; \
            echo "Generated gRPC code for $$prefix"; \
        fi; \
    done
	


# 运行Flutter应用，用于开发和测试
run_frontend:
	@echo "Running Flutter app..."
	cd frontend && flutter pub get
	cd frontend && flutter run

# 运行Go后端服务
run_backend: stop_backend check_ports
	@mkdir -p backend/pids backend/logs
	@echo "Running Go backend api-gateway service"
	@(cd backend/go-services/api-gateway/ && nohup go run main.go > $(CURDIR)/backend/logs/api-gateway.log 2>&1 & echo $$! > $(CURDIR)/backend/pids/api-gateway.pid)

	@echo "Running Go backend auth service"
	@(cd backend/go-services/auth/ && nohup go run main.go > $(CURDIR)/backend/logs/auth.log 2>&1 & echo $$! > $(CURDIR)/backend/pids/auth.pid)

	@echo "Running Go backend message service"
	@(cd backend/go-services/message/ && nohup go run main.go > $(CURDIR)/backend/logs/message.log 2>&1 & echo $$! > $(CURDIR)/backend/pids/message.pid)

	@echo "Running Go backend poll service"
	@(cd backend/go-services/poll/ && nohup go run main.go > $(CURDIR)/backend/logs/poll.log 2>&1 & echo $$! > $(CURDIR)/backend/pids/poll.pid)

	@echo "Running Go backend recommendation service"
	@(cd backend/go-services/recommendation/ && nohup go run main.go > $(CURDIR)/backend/logs/recommendation.log 2>&1 & echo $$! > $(CURDIR)/backend/pids/recommendation.pid)

	@echo "Running Java backend user service"
	@(cd backend/java-services/user/ && nohup mvn spring-boot:run > $(CURDIR)/backend/logs/user.log 2>&1 & echo $$! > $(CURDIR)/backend/pids/user.pid)

# 停止Go后端服务
stop_backend:
	@echo "Stopping backend services..."
	-@kill `cat backend/pids/api-gateway.pid` 2>/dev/null || true
	-@kill `cat backend/pids/auth.pid` 2>/dev/null || true
	-@kill `cat backend/pids/message.pid` 2>/dev/null || true
	-@kill `cat backend/pids/poll.pid` 2>/dev/null || true
	-@kill `cat backend/pids/recommendation.pid` 2>/dev/null || true
	-@kill `cat backend/pids/user.pid` 2>/dev/null || true
	@rm -f backend/pids/*.pid

# 检查端口占用情况
check_ports:
	@for port in 8080 8081 8082 8083 8084 8085; do \
		if lsof -i:$$port > /dev/null; then \
			echo "Port $$port is in use by:"; \
			lsof -i:$$port; \
			read -p "Do you want to kill the process using port $$port? (y/n) " choice; \
			if [ "$$choice" = "y" ]; then \
				fuser -k $$port/tcp; \
				echo "Killed the process using port $$port"; \
			else \
				echo "Port $$port is in use. Aborting..."; \
				exit 1; \
			fi; \
		fi; \
	done

# 构建整个项目，包括前端和后端
build: flutter_build go_build

# 构建Flutter前端应用,用于发布应用
flutter_build:
	@echo "Building Flutter app for all platforms..."
	cd frontend && flutter build apk
	cd frontend && flutter build ios
	cd frontend && flutter build web
	cd frontend && flutter build linux
	cd frontend && flutter build macos
	cd frontend && flutter build windows

# 构建Go后端服务
go_build:
	@echo "Building Go backend services..."
	@echo "Building Go api-gateway service"
	(cd backend/go-services/api-gateway && go build -o $(CURDIR)/build/api-gateway)
	@echo "Building Go auth service"
	(cd backend/go-services/auth && go build -o $(CURDIR)/build/auth)
	@echo "Building Go message service"
	(cd backend/go-services/message && go build -o $(CURDIR)/build/message)
	@echo "Building Go poll service"
	(cd backend/go-services/poll && go build -o $(CURDIR)/build/poll)
	@echo "Building Go recommendation service"
	(cd backend/go-services/recommendation && go build -o $(CURDIR)/build/recommendation)

	@echo "Building Java backend user service"
	(cd backend/java-services/user && mvn package)

# 清理构建文件
clean:
	@echo "Cleaning build artifacts..."
	cd frontend && flutter clean
	find . -type f -name '*.lock' -delete
	find . -type f -name '*.log' -delete
	rm -rf build/backend_service
	rm -rf build/api-gateway build/auth build/message build/poll build/recommendation
	rm -rf backend/java-services/user/target

# 部署项目
deploy: kubernetes_deploy monitoring_deploy

# 部署到Kubernetes
kubernetes_deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f deployments/kubernetes/

# 部署监控配置
monitoring_deploy:
	@echo "Deploying monitoring stack..."
	kubectl apply -f deployments/monitoring/
