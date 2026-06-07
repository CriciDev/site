FROM golang:1.26.1-alpine AS build

WORKDIR /src

COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o /out/site .

FROM scratch

WORKDIR /app

COPY --from=build /out/site /app/site
COPY --from=build /src/assets /app/assets
COPY --from=build /src/db /app/db

EXPOSE 8080

ENTRYPOINT ["/app/site"]
