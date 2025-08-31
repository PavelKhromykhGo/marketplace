FROM golang:1.24-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /cmd/marketplace
RUN go run github.com/swaggo/swag/cmd/swag@latest init \
    --generalInfo ./main.go \
    --dir .,/app/internal \
    --output app/docs

WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app/marketplace ./cmd/marketplace

FROM alpine:3.19
WORKDIR /app
COPY --from=builder app/marketplace app/marketplace
COPY migrations app/migrations
COPY --from=builder app/docs app/docs

EXPOSE 8080
ENV DATABASE_URL="host=postgres port=5432 user=postgres password=postgres dbname=marketplace sslmode=disable"
ENV MIGRATIONS_DIR="/app/migrations"
ENTRYPOINT ["/app/marketplace"]
