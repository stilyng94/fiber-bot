lint:
	@docker run --rm -v $(pwd):/app -v ~/.cache/golangci-lint/v1.55.1:/root/.cache -w /app golangci/golangci-lint:v1.55.1 golangci-lint run -v

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
	@go build -o bin/server cmd/api/main.go
	@echo "build done!!"

run: build
	@./bin/server

frontend-watch:
	@cd cmd/frontend && pnpm watch

frontend-build:
	@echo "building frontend ...."; \
  cd cmd/frontend && pnpm build
	@echo "frontend build done ...."

tunnel:
	@ngrok http --host-header="localhost:${port}" ${port}


stop_tunnel:
	@echo "Stopping tunnel..."
	@-pkill -SIGTERM -f "tunnel -port=${port}"
	@echo "Stopped tunnel"

migrate-build:
	@echo "building migrate ...."
	@go build -o bin/migrate cmd/migrate/main.go
	@echo "build done!!"

migrate-run:
	@./bin/migrate

docker_run:
	@docker run --rm --name fiber-bot -p 5005:5005 -v $(pwd)/.env:/app/.env:ro escobar0216/fiber-bot:latest ./migrate db help
