.PHONY: proto clean run

# Generate protobuf code
proto:
	@mkdir -p gen
	protoc -I./proto -I$$HOME/.proto \
		--go_out=./gen --go_opt=paths=source_relative \
		--go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./gen --grpc-gateway_opt=paths=source_relative \
		proto/*.proto
	@echo "✓ Proto files generated successfully"

# Clean generated files
clean:
	rm -rf gen/*.pb.go gen/*.pb.gw.go
	@echo "✓ Generated files cleaned"

# Run the server
run:
	go run main.go
