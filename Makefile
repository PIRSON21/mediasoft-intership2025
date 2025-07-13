GO_FILES := $(shell find . -type f -name '*.go')
CONFIG_FILE := ./configs/make.env

include ${CONFIG_FILE}
export $(shell sed 's/=.*//' ${CONFIG_FILE})

run: build
	@cd build ; ./app

build: ${GO_FILES} ${CONFIG_FILE}
	@echo "\n🛠️  Building project..."
	@mkdir ./build/static -p
	@go build -o ./build/app ./cmd/intership/. || (echo "❌ Build failed!" && exit 1)
	@echo "✔️  Success builded!\n"

create_migration:
	@if [ -z "$(MIG_NAME)" ]; then \
		echo "❌ Не указано NAME: make create_migration MIG_NAME=название миграции"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir db/migrations/ -seq $(MIG_NAME)

up_migration:
	@migrate -path db/migrations/ -database $(DB_URL) up $(COUNT)

down_migration:
	@migrate -path db/migrations/ -database $(DB_URL) down $(COUNT)

swagger:
	@swagger generate spec -o ./api/swagger.json --scan-models

docker-up:
	@if docker compose --env-file configs/docker.env --file deployments/docker-compose.yml up --wait -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling Docker Compose Up"; \
		exit 1; \
	fi


docker-down:
	@if docker compose --file deployments/docker-compose.yml down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling Docker Compose Down"; \
		exit 1; \
	fi

.PHONY: run build create_migration up_migration down_migration swagger docker-up docker-down
