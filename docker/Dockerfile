FROM golang:1.22 AS builder

ENV GOTOOLCHAIN=auto
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migrate migrations/migration_001.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/seed seeds/seed_users.go

FROM alpine:3.20
WORKDIR /app

ENV PORT=8080
ENV DB_PATH=/data/data.db
ENV RUN_MIGRATIONS=1
ENV RUN_SEED=0

RUN addgroup -S app && adduser -S app -G app
RUN mkdir -p /data && chown -R app:app /data /app

COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/migrate /app/migrate
COPY --from=builder /app/bin/seed /app/seed
COPY docker/docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

USER app
EXPOSE 8080
VOLUME ["/data"]

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/app/server"]
