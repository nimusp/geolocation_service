# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -v -o gateway cmd/gateway/main.go

EXPOSE 8888

CMD [ "/app/gateway" ]