FROM golang:latest AS builder
WORKDIR /opt
COPY . .
RUN cp .env.example .env \
 && apt update \
 && GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -mod vendor -o /opt/wb-app ./cmd/main.go

FROM alpine:latest

WORKDIR /opt
COPY --from=builder ["/opt/*.env", "/opt/wb-app", "/opt/"]

ENTRYPOINT ["./wb-app"]
