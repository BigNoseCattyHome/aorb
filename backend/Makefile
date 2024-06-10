IDL_PATH=$(shell pwd)/idl
SERVICES:=gateway auth comment question recommend user vote

.PHONY: proto
proto:
	@echo "Creating grpc file..."
	@for file in $(IDL_PATH)/*.proto; do	\
  		if [ -f "$$file" ]; then	\
  		  prefix=$$(basename "$$file" .proto);	\
  		  mkdir -p ./rpc/"$${prefix}";	\
  		  echo "Created dictionary for $$prefix";	\
  		  protoc --proto_path=./idl/ --go_out=./rpc/$$prefix --go_opt=paths=source_relative --go-grpc_out=./rpc/$$prefix	\
  		  	--go-grpc_opt=paths=source_relative ./idl/$$prefix.proto;	\
  		fi;	\
  	done
