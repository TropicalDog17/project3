# syntax=docker/dockerfile:1
FROM golang:1.21 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o gset gset.go

FROM alpine:3.14
ENV RESOLUTION 1366x768x24
ENV TIME_LIMIT 20
ENV NODE_COUNT 3
ENV RATE 100
RUN  apk update \
    && apk upgrade \
    && apk add ca-certificates \
    && update-ca-certificates \
    && apk add --update coreutils && rm -rf /var/cache/apk/*   \ 
    && apk add --update openjdk11 tzdata curl unzip bash \
    && apk add --no-cache nss \
    && apk add --no-cache git \
    && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build /app/gset gset
COPY maelstrom maelstrom
CMD java -Djava.awt.headless=true -jar maelstrom/lib/maelstrom.jar test -w g-set --bin gset --node-count ${NODE_COUNT} --rate ${RATE} --time-limit ${TIME_LIMIT} --nemesis partition
