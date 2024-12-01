.PHONY: proto clean server client build

# Build targets
proto:
	mkdir -p proto/generated
	protoc --go_out=./proto/generated --go_opt=paths=source_relative --go-grpc_out=./proto/generated/ --go-grpc_opt=paths=source_relative proto/document_service.proto

build:
	mkdir -p bin
	go build -o bin/server cmd/server/*.go
	go build -o bin/client cmd/client/*.go

# Clean targets
clean:
	rm -rf proto/generated
	rm -rf bin/

# Server target
server:
	go run cmd/server/*.go -port 50051

# Interactive client target
client:
	go run cmd/client/*.go -server localhost:50051