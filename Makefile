 # Regenerate gRPC code 
proto:
	protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative  ./api/api.proto

.PHONY: proto