# 设置默认目标不是文件
.PHONY: all clean test build deploy run_backend run_frontend

# 默认执行的命令
all: build

# 运行Flutter应用，用于开发和测试
run_frontend:
	@echo "Running Flutter app..."
	cd frontend && flutter pub get
	cd frontend && flutter run

# 运行Go后端服务
run_backend:
	@echo "Running Go backend service..."
	cd backend/build && ./backend_service


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
	@echo "Building Go backend service..."
	cd backend/api && go build -o ../build/backend_service

# 清理构建文件
clean:
	@echo "Cleaning build artifacts..."
	flutter clean
	find . -type f -name '*.lock' -delete
	find . -type f -name '*.log' -delete
	rm -f build/backend_service

# 运行所有测试
test: flutter_test go_test

# 运行Flutter测试
flutter_test:
	@echo "Running Flutter tests..."
	flutter test

# 运行Go测试
go_test:
	@echo "Running Go tests..."
	cd backend/api && go test ./...

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
