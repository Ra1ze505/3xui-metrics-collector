FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /exporter ./cmd/exporter

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /exporter /app/exporter
EXPOSE 2112
ENTRYPOINT ["/app/exporter"]
CMD ["-config", "/app/config.yaml"]
