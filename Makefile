.PHONY: all clean build deploy run_backend run_frontend proto stop

# 默认执行的命令
all: proto build

PROTO_PATH=./idl
GOOGLE_PROTO_PATH=/usr/local/include
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
			protoc -I$(PROTO_PATH) -I$(GOOGLE_PROTO_PATH) \
				--go_out=$(GO_OUT_PATH)/$$prefix --go_opt=paths=source_relative \
				--go-grpc_out=$(GO_OUT_PATH)/$$prefix --go-grpc_opt=paths=source_relative \
				--dart_out=grpc:$(OUTPUT_DART_PATH) \
				--plugin=protoc-gen-dart=$(PROTOC_GEN_DART) \
				$$file; \
			echo "Generated gRPC code for $$prefix"; \
		fi; \
	done
	@protoc -I$(GOOGLE_PROTO_PATH) \
		--dart_out=grpc:$(OUTPUT_DART_PATH) \
		--plugin=protoc-gen-dart=$(PROTOC_GEN_DART) \
		google/protobuf/timestamp.proto
	@echo "Generated Dart code for Google's timestamp.proto"

# 构建整个项目，包括前端和后端
build:  flutter_build go_build

# 构建Flutter前端应用,用于发布应用
flutter_build:
	@echo "Building Flutter app for all platforms..."
	cd frontend && flutter pub get
	# cd frontend && flutter build apk
	# cd frontend && flutter build ios
	# cd frontend && flutter build web
	cd frontend && flutter build linux
	# cd frontend && flutter build macos
	# cd frontend && flutter build windows

# 构建Go后端服务
go_build:
	@echo "Building Go backend service..."
	bash ./scripts/build_all_backend.sh

# 清理构建文件，需要重新build
clean:
	@echo "Cleaning build artifacts..."
	cd frontend && flutter clean
	find . -type f -name '*.lock' -delete
	find . -type f -name '*.log' -delete
	rm -rf build/*

run: run_frontend run_backend

# 运行Flutter应用，run之前需要`make build`
run_frontend:
	@echo "Running Flutter app..."
	cd frontend && flutter run

# 运行Go后端服务
run_backend:
	@echo "Running Go backend service..."
	bash ./scripts/run_all_backend.sh

# 停止运行后端服务
stop:
	@echo "Stopping all services..."
	bash ./scripts/stop_all_backend.sh

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
