dev-server:
	nodemon --exec go run main.go --signal SIGTERM

pre-commit:
	pre-commit run --all-files

PHONY: start-dev pre-commit
