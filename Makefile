dev-server:
	nodemon --exec go run main.go --signal SIGTERM

pre-commit:
	pre-commit run --all-files

swagger:
	swag init --dir cmd/api/,internal/schemas/,internal/auth/handler

swagger_format:
	swag fmt --dir cmd/api/,internal/schemas/,internal/auth/handler

PHONY: start-dev pre-commit swagger swagger_format
