FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm64
RUN go build -o /app/gateway .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/gateway /app/gateway

EXPOSE 8000
ENTRYPOINT ["/app/gateway"]