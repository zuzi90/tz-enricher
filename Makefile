GREEN 		= \033[0;32m
YELLOW 		= \033[0;33m
NC 			= \033[0m

up:
	docker compose up -d

down:
	docker compose down

test: up
	@echo "\n${GREEN}Running unit-tests${NC}"
	go test ./...
	make down

lint:
	@echo "\n${GREEN}Linting Golang code with golangci${NC}"
	gofumpt -w .
	go mod tidy
	golangci-lint run ./...

run:
	@echo "\n${GREEN}Run the application${NC}"
	go run cmd/main/main.go

swag:
	@echo "\n${GREEN}Generate Swagger documentation${NC}"
	swag init -g ./cmd/main/main.go