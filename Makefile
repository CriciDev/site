.PHONY: help run run-docker run-dev run-dev-docker down templ test build clean

help:
	@printf '%s\n' \
		'Targets:' \
		'  make run         - roda a app sem hot reload no host' \
		'  make run-docker  - sobe a app empacotada via docker compose' \
		'  make run-dev     - hot reload no host com air + templ' \
		'  make run-dev-docker - hot reload no Docker com air + templ' \
		'  make down        - derruba containers do profile dev' \
		'  make templ       - gera arquivos templ' \
		'  make test        - roda testes Go' \
		'  make build       - compila o binario em ./tmp/site' \
		'  make clean       - remove artefatos locais de build'

down:
	docker compose --profile dev down

templ:
	templ generate

test:
	go test ./...

build:
	@mkdir -p tmp
	go build -o ./tmp/site ./cmd/site

run:
	go run ./cmd/site

run-docker:
	docker compose up -d --build

run-dev:
	./scripts/dev

run-dev-docker:
	docker compose --profile dev up --build app-dev

clean:
	rm -rf ./tmp ./build-errors.log
