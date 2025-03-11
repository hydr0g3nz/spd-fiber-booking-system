#!/bin/bash

# This script generates Swagger documentation using swaggo/swag
# Ensure you have swag installed: go install github.com/swaggo/swag/cmd/swag@latest

echo "Generating Swagger documentation..."

# Remove existing docs
rm -rf ./docs/docs.go
rm -rf ./docs/swagger.json
rm -rf ./docs/swagger.yaml

# Generate new docs
swag init -g cmd/main.go -o docs

echo "Swagger documentation generated successfully!"
echo "Run the application and visit http://localhost:3000/swagger/ to see the documentation."