FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD [ "go", "run", "main.go", "-env=pgsql"]