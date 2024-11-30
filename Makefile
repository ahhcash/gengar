.PHONY: proto

proto:
	mkdir -p proto/generated
	protoc --go_out=./proto/generated --go_opt=paths=source_relative --go-grpc_out=./proto/generated/ --go-grpc_opt=paths=source_relative proto/document_service.proto

clean:
	rm -rf proto/generated