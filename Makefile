gen-proto:
	# Generate protobuf and gRPC files
	protoc --go_out=proto-gen --go-grpc_out=proto-gen pb/auth.proto

image:
	# Build Docker image
	docker build -t api-gateway:latest .

container-run:
	# Run Docker container
	docker run -p 8080:8080 api-gateway:latest
