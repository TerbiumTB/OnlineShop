FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder
LABEL authors="terbium"

WORKDIR /app

COPY go.mod ./

RUN go mod tidy

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/payments .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/payments /app/payments

EXPOSE 8080
ENTRYPOINT ["/app/payments"]