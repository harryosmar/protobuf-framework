.PHONY: proto clean run build swagger


# ==== Configuration ====
# Your Go module path (used if MODE=import)
PROJECT_MOD       ?= github.com/harryosmar/protobuf-go

# Generate protobuf code
proto:
	@mkdir -p gen
	@export PATH="/usr/local/go/bin:$$PATH" && \
	protoc -I./proto -I./third_party -I$$HOME/.proto \
		--go_out=module=$(PROJECT_MOD),paths=import:. \
		--go-grpc_out=module=$(PROJECT_MOD),paths=import:. \
		--grpc-gateway_out=module=$(PROJECT_MOD),paths=import:. \
		--gorm_out=module=$(PROJECT_MOD),paths=import:. \
		--validate_out=lang=go,module=$(PROJECT_MOD),paths=import:. \
		--go-scaffold_out=base=$(PROJECT_MOD),paths=source_relative:. \
		proto/*.proto
	@echo "✓ Proto files generated successfully with validation and GORM models"

# Generate Swagger/OpenAPI documentation
swagger:
	@mkdir -p docs
	protoc -I./proto -I./third_party -I$$HOME/.proto \
		--openapiv2_out=./docs \
		--openapiv2_opt=logtostderr=true \
		--openapiv2_opt=allow_merge=true \
		--openapiv2_opt=merge_file_name=api \
		proto/*.proto
	@echo "✓ Swagger documentation generated in docs/api.swagger.json"

# Clean generated files
clean:
	rm -rf gen/**/*.pb.go gen/**/*.pb.gw.go gen/**/*.pb.validate.go gen/**/*.gorm.go docs/*.swagger.json
	@echo "✓ Generated files cleaned"

# Build the application
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
	@echo "✓ Application built successfully"

# Run the server
run:
	go run main.go
