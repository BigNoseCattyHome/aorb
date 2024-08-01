.PHONY: all clean build deploy run stop proto frontend backend

# 变量定义
PROTO_PATH := ./idl
GOOGLE_PROTO_PATH := /usr/local/include
OUTPUT_DART_PATH := ./frontend/lib/generated/
GO_OUT_PATH := ./backend/rpc
PROTOC_GEN_DART := $(shell which protoc-gen-dart)

# 默认目标
all: proto build

# 生成 protobuf 文件
proto:
	@echo "Creating golang and dart grpc files..."
	@for file in $(PROTO_PATH)/*.proto; do \
		if [ -f "$$file" ]; then \
			prefix=$$(basename "$$file" .proto); \
			mkdir -p $(GO_OUT_PATH)/"$${prefix}" $(OUTPUT_DART_PATH); \
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

# 构建项目
build:
	@case "$(filter-out $@,$(MAKECMDGOALS))" in \
		all) make flutter-build go-build ;; \
		frontend) make flutter-build ;; \
		backend) make go-build ;; \
		*) echo "Usage: make build [all|frontend|backend]" ;; \
	esac

flutter-build:
	@echo "Building Flutter app for all platforms..."
	cd frontend && flutter build linux
	cd frontend && flutter build apk
	cd frontend && flutter build ios
	# Uncomment below lines as needed
	# cd frontend && flutter build web
	# cd frontend && flutter build macos
	# cd frontend && flutter build windows

go-build:
	@echo "Building Go backend service..."
	bash ./scripts/build_all_backend.sh

# 清理构建文件
clean:
	@echo "Cleaning build artifacts..."
	cd frontend && flutter clean
	find . -type f \( -name '*.lock' -o -name '*.log' \) -delete
	rm -rf build/*

# 运行服务
run:
	@case "$(filter-out $@,$(MAKECMDGOALS))" in \
		backend) make run-backend ;; \
		frontend) make run-frontend ;; \
		*) echo "Usage: make run [backend|frontend]" ;; \
	esac

run-backend:
	@echo "Running Go backend service..."
	bash ./scripts/check_and_stop_ports.sh
	bash ./scripts/run_backend.sh

run-frontend:
	@echo "Running Flutter app..."
	cd frontend && flutter run

# 停止服务
stop:
	@echo "Stopping all services..."
	bash ./scripts/check_and_stop_ports.sh

# 部署项目
deploy: kubernetes-deploy monitoring-deploy

kubernetes-deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f deployments/kubernetes/

monitoring-deploy:
	@echo "Deploying monitoring stack..."
	kubectl apply -f deployments/monitoring/

%:
	@: