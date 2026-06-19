.PHONY: help volume up down logs restart build run dev dev-docker templ test clean

ifeq ($(OS),Windows_NT)
DEV_CMD = powershell -NoProfile -ExecutionPolicy Bypass -File scripts/dev.ps1
else
DEV_CMD = ./scripts/dev
endif

help:
	@printf '%s\n' \
		'Targets:' \
		'  make up          - sobe a app empacotada com Docker' \
		'  make down        - derruba os containers do compose' \
		'  make logs        - acompanha logs do compose' \
		'  make restart     - reinicia o ambiente Docker' \
		'  make run         - roda a app no host sem hot reload' \
		'  make dev         - hot reload no host com air + templ' \
		'  make dev-docker  - hot reload no Docker com air + templ' \
		'  make templ       - gera arquivos templ' \
		'  make test        - roda testes Go' \
		'  make build       - compila o binario em ./tmp/site' \
		'  make clean       - remove artefatos locais de build'

volume:
	docker volume create sqlite_data

up: volume
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f

restart: down up

build:
	@mkdir -p tmp
	go build -o ./tmp/site ./cmd/site

run:
	go run ./cmd/site

dev:
	@$(DEV_CMD)

dev-docker: volume
	docker compose --profile dev up --build app-dev

templ:
	templ generate

test:
	go test ./...

clean:
	rm -rf ./tmp ./build-errors.log
