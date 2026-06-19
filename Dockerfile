FROM golang:1.26.1-alpine AS dev

RUN apk add --no-cache bash git

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/cosmtrek/air@v1.51.0
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020

FROM golang:1.26.1-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/site ./cmd/site

FROM scratch

WORKDIR /app

COPY --from=build /out/site /app/site
COPY --from=build /src/assets /app/assets
COPY --from=build /src/db /app/db

EXPOSE 8080

ENTRYPOINT ["/app/site"]
