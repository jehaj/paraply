FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /will-rain

FROM ubuntu:latest

COPY --from=builder /will-rain /will-rain
USER nonroot:nonroot
ENTRYPOINT ["/will-rain"]
