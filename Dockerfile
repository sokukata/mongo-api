
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY convert/ ./convert/
COPY mongoclient/ ./mongoclient/

RUN go build -o /mongo-api

EXPOSE 8080

ENTRYPOINT [ "/mongo-api" ]