 # Regenerate gRPC code 
proto:
	protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative  ./api/api.proto

docker:
# docker build -t image-name:tag dockerfile-path
	docker build -t ledger:v0.1 .

.PHONY: proto docker