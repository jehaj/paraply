FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /paraply

FROM ubuntu:latest

COPY --from=builder /paraply /paraply
USER nonroot:nonroot
ENTRYPOINT ["/paraply"]
