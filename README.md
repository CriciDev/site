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

### Com Make

Para não decorar os comandos:

```bash
make help
```

Os atalhos principais são:

- `make run`: roda no host sem hot reload
- `make run-docker`: sobe a versão empacotada com Docker
- `make run-dev`: hot reload no host
- `make run-dev-docker`: hot reload no Docker
- `make test`: testes Go
- `make build`: build local em `tmp/site`

### Com Go

```bash
go run ./cmd/site
```

Por padrão a aplicação sobe em `:8080` e grava o banco em `data/cricidev.db`.

Se quiser mudar o caminho do banco:

```bash
DB_PATH=/outro/caminho/cricidev.db go run ./cmd/site
```

### Com hot reload no modo dev

O projeto agora tem um runner de desenvolvimento com:

- `air` para rebuild/restart do binário Go
- `templ` para regenerar templates e fazer live reload no navegador

Rode:

```bash
make run-dev
```

Abra no navegador:

```text
http://localhost:7331
```

Esse endereço é o proxy do `templ`, que recarrega a página quando você altera:

- arquivos `.go`
- arquivos `.templ`
- assets em `assets/` como `.css`, `.js` e imagens

O servidor Go continua rodando internamente em `http://localhost:8080`.

Se preferir, o mesmo fluxo também pode rodar dentro do Docker.

### Com hot reload via Docker

Essa opção sobe um container de desenvolvimento com:

- `air` dentro do container
- `templ` dentro do container
- bind mount do código local
- SQLite salvo em `data/cricidev.db` no próprio repositório

Rode:

```bash
make run-dev-docker
```

Abra no navegador:

```text
http://localhost:7331
```

Notas:

- não rode `app` e `app-dev` ao mesmo tempo, porque ambos usam as portas `8080` e `7331`
- no modo Docker dev, o banco fica em `data/cricidev.db` no host via bind mount
- se quiser derrubar depois, use `make down`

### Com Docker Compose

O `docker compose` usa um volume externo chamado `sqlite_data`.

Se ele ainda não existir:

```bash
docker volume create sqlite_data
```

Depois:

```bash
make run-docker
```

Nessa forma, o banco fica em `/data/cricidev.db` dentro do container e é persistido pelo volume.

## Variáveis de ambiente

- `ADDR`: endereço HTTP da aplicação. Padrão `:8080`
- `DB_PATH`: caminho do arquivo SQLite. Padrão `data/cricidev.db` fora do Docker e `/data/cricidev.db` no `docker compose`

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

No fluxo de desenvolvimento com `./scripts/dev`, isso acontece automaticamente.

## Build

O `Dockerfile` compila o binário a partir de `./cmd/site`.
