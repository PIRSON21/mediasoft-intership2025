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
	@swagger generate spec -o ./api/swagger.yaml --scan-models
