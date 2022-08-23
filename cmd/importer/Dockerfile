# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -v -o importer cmd/importer/main.go

CMD [ "/app/importer" ]