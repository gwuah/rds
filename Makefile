# include ./.env

BRANCH ?= main
BUILD_N ?= 0

build:
	go build -ldflags="-X 'main.Version=1.0.0.$(BUILD_N)-$(BRANCH)'" -o ./bin/rds ./cmd/rds

run: build
	@DATABASE_URL=${DATABASE_URL} REDIS_URL=${REDIS_URL} SLACK_URL=${SLACK_URL} bin/rds --port 8091 --address 0.0.0.0

migration:
	@goose --dir internal/db/migrations create $(name) sql

psql:
	@psql ${DATABASE_URL}

api-gen:
	cd api && buf lint
	cd api && buf generate