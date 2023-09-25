lint:
	@docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.52.2 golangci-lint run -v

create_model:
	@echo "Creating new model..."
	@go run -mod=mod entgo.io/ent/cmd/ent new ${model}

dev:
	@echo "Starting api in dev mode..."
	@air
	@echo "api running!"

generate:
	@echo "Genrating...."
	@go generate ./...
	@echo "Genrating done!!"

build: frontend-build
	@echo "building server ...."
	@go build -o ./bin/server
	@echo "build done!!"

run: build
	@./bin/server

frontend-watch:
	@cd www/fiber-bot && pnpm watch

frontend-build:
	@echo "building frontend ...."; \
  cd www/fiber-bot && pnpm build
	@echo "frontend build done ...."
