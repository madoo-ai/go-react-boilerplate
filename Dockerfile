# syntax=docker/dockerfile:1

FROM node:24-bookworm-slim AS web
WORKDIR /src
COPY web/package.json web/package-lock.json web/
RUN --mount=type=cache,target=/root/.npm cd web && npm ci --no-audit --no-fund
COPY web web
RUN mkdir -p internal/handler/webui && cd web && npm run build

FROM golang:1.24-bookworm AS gobuild
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
COPY --from=web /src/internal/handler/webui/dist ./internal/handler/webui/dist
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/server ./cmd/server \
 && CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/seed ./cmd/seed \
 && CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.22 AS runner
RUN apk add --no-cache ca-certificates \
 && adduser -D -u 10001 app
WORKDIR /app
COPY --from=gobuild /out/server /app/server
COPY --from=gobuild /out/seed /app/seed
COPY --from=gobuild /out/migrate /app/migrate
USER app
ENV APP_ENV=production
ENV LISTEN=0.0.0.0:3000
EXPOSE 3000
ENTRYPOINT ["/app/server"]
