FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/go-service ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /out/go-service /app/go-service

EXPOSE 8080

CMD ["/app/go-service"]
