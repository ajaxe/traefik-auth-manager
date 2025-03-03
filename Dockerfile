FROM golang:alpine AS builder

ARG InstallFolder=/go/src/github.com/ajaxe/traefik-auth-manager

RUN mkdir -p $InstallFolder

WORKDIR $InstallFolder

COPY . .

RUN GOARCH=wasm GOOS=js go build -o ./tmp/web/app.wasm ./cmd/webapp \
    && go build -o ./tmp/server ./cmd/webapp/ \
    && cp -a ./web/* ./tmp/web/

FROM alpine:latest AS runner

ARG InstallFolder=/go/src/github.com/ajaxe/traefik-auth-manager \
    Port=8000

RUN apk add --no-cache curl

RUN mkdir -p /home/app/
WORKDIR /home/app/

COPY --from=builder "${InstallFolder}/tmp/." .

ENV APP_ENV=production \
    APP_SERVER_PORT=$Port

EXPOSE $Port

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl --fail "http://localhost:8000/healthcheck" || exit

CMD ["./server"]
