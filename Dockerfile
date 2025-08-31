FROM golang:1.24-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/marketplace ./cmd/marketplace

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /bin/marketplace app/bin/marketplace
COPY migrations app/migrations

EXPOSE 8080
ENV DATABASE_URL="host=postgres port=5432 user=postgres password=postgres dbname=marketplace sslmode=disable"
ENV MIGRATIONS_DIR="/app/migrations"
ENTRYPOINT ["/app/marketplace"]
