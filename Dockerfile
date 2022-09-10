# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY . .
RUN go mod download


RUN go build -o /app/bin/kaomoji-db-api

EXPOSE 3000

CMD [ "/app/bin/kaomoji-db-api" ]
