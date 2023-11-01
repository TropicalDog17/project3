# syntax=docker/dockerfile:1
FROM golang:1.21 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o gset gset.go

FROM ubuntu:22.04
ENV RESOLUTION 1366x768x24
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive \
    apt-get install openjdk-11-jdk -y && apt-get install git gnuplot -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /app/gset gset
COPY maelstrom maelstrom
CMD ["java", "-Djava.awt.headless=true", "-jar","maelstrom/lib/maelstrom.jar","test", "-w", "g-set",  "--bin", "gset", "--node-count", "3", "--rate", "100", "--time-limit" ,"20" ,"--nemesis", "partition"]

