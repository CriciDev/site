# Criciúma Devs Site

Site oficial da comunidade Criciúma Devs.

O projeto usa:

- Go
- `templ`
- HTMX
- SQLite
- Docker

## O que tem aqui

- landing page da comunidade
- coleta simples de analytics em SQLite
- assets estáticos em `assets/`

## Estrutura

```text
cmd/site         entrypoint da aplicação
internal/app     bootstrap da aplicação
internal/server  rotas HTTP
internal/analytics  persistência e eventos
pages/home       template da home
assets/          CSS, JS e imagens
db/              schema do SQLite
```

## Requisitos

- Go 1.26+
- Docker, se for rodar com `docker compose`

## Rodando localmente

### Com Go

```bash
go run ./cmd/site
```

Por padrão a aplicação sobe em `:8080` e grava o banco em `/data/cricidev.db`.

### Com Docker Compose

O `docker compose` usa um volume externo chamado `sqlite_data`.

Se ele ainda não existir:

```bash
docker volume create sqlite_data
```

Depois:

```bash
docker compose up -d --build
```

## Variáveis de ambiente

- `ADDR`: endereço HTTP da aplicação. Padrão `:8080`
- `DB_PATH`: caminho do arquivo SQLite. Padrão `/data/cricidev.db`

## Analytics

O site envia eventos do navegador para `POST /api/events`.

Eventos registrados:

- `pageview`
- `session_start`
- `session_end`
- `click`

Os eventos são salvos nas tabelas `events` e `sessions`.

## HTMX

O `htmx` já é carregado na home via CDN e fica disponível para uso em qualquer tela nova.

## Templates

O HTML da home é gerado com `templ`.

Depois de alterar arquivos `.templ`, rode:

```bash
templ generate
```

## Build

O `Dockerfile` compila o binário a partir de `./cmd/site`.
