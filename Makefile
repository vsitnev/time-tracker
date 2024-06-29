include .env
export

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ### remove docker volume
	docker volume rm sync-manager-db
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' -c --ignored | xargs hadolint
.PHONY: linter-hadolint

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations 'time-tracker'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### migration down
	echo "y" | migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' down
.PHONY: migrate-down

swag: ### generate swagger docs
	swag init -g ./cmd/app/main.go
