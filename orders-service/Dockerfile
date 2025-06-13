FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder
LABEL authors="terbium"

WORKDIR /app

COPY go.mod ./

RUN go mod tidy

RUN go mod download

COPY . .


RUN CGO_ENABLED=0 go build -o /app/orders .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/orders /app/orders

EXPOSE 8080
ENTRYPOINT ["/app/orders"]