FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o url-shortener ./cmd

FROM golang:1.23-alpine AS runner

WORKDIR /app

COPY --from=builder /app/url-shortener /app/url-shortener

COPY ../config/config.yaml /app/config/config.yaml

CMD ["/app/url-shortener"]
