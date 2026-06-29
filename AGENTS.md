# AGENTS.md

## Contexto
- Site oficial da comunidade Criciúma Devs.
- Stack: Go, `templ`, HTMX, SQLite e Docker.
- A home é renderizada com `templ` e os eventos de analytics são salvos em SQLite.

## Estrutura principal
- `cmd/site`: entrypoint da aplicação.
- `internal/app`: bootstrap, banco e servidor HTTP.
- `internal/server`: rotas e handlers.
- `internal/analytics`: persistência e modelo de eventos.
- `pages/home`: template da home.
- `assets/`: CSS, JS e imagens.
- `db/schema.sql`: schema do SQLite.

## Regras de edição
- Edite `pages/home/page.templ`, nunca `page_templ.go` manualmente.
- Depois de alterar `.templ`, rode `templ generate`.
- Prefira mudanças pequenas e seguras.
- Mantenha nomes de classes CSS em camelCase.
- Evite introduzir novos nomes no padrão BEM.
- Não mexa em arquivos gerados quando a mudança puder ser feita na fonte.

## Estilo
- Go: siga o estilo idiomático padrão.
- CSS: camelCase para classes.
- JS: funções pequenas, sem adicionar frameworks extras.
- HTML/templ: direto e legível, sem abstração desnecessária.

## Validação
- Use `make check` antes de concluir mudanças.
- Se alterar templates, rode `templ generate` antes de validar.
- Se mexer em analytics ou server, confirme com `go test ./...` e `make check`.

## Comandos úteis
- `make dev`: hot reload no host.
- `make dev-docker`: hot reload no Docker.
- `make up`: sobe a versão Docker.
- `make test`: roda testes.
- `make build`: compila em `tmp/site`.
- `make check`: fmt, vet, test e build.

## Banco e arquivos gerados
- Banco local padrão: `data/cricidev.db`.
- Arquivos em `data/*.db*` e `tmp/` são gerados e não devem ser versionados.
- O schema do banco vem de `db/schema.sql`.

## Boas práticas
- Leia o código existente antes de propor refatorações.
- Evite renomeações amplas sem necessidade.
- Preserve o comportamento visual e os fluxos existentes.
- Se uma mudança exigir quebra, sinalize antes.
