.PHONY: proto clean run build swagger

# Generate protobuf code
proto:
	@mkdir -p gen
	protoc -I./proto -I$$HOME/.proto \
		--go_out=. --go_opt=paths=import \
		--go-grpc_out=. --go-grpc_opt=paths=import \
		--grpc-gateway_out=. --grpc-gateway_opt=paths=import \
		proto/*.proto
	@echo "✓ Proto files generated successfully"

# Generate Swagger/OpenAPI documentation
swagger:
	@mkdir -p docs
	protoc -I./proto -I$$HOME/.proto \
		--openapiv2_out=./docs \
		--openapiv2_opt=logtostderr=true \
		--openapiv2_opt=allow_merge=true \
		--openapiv2_opt=merge_file_name=api \
		proto/*.proto
	@echo "✓ Swagger documentation generated in docs/api.swagger.json"

# Clean generated files
clean:
	rm -rf gen/*.pb.go gen/*.pb.gw.go docs/*.swagger.json
	@echo "✓ Generated files cleaned"

# Build the application
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
	@echo "✓ Application built successfully"

# Run the server
run:
	go run main.go
