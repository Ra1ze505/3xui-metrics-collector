FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o 3xui-metrics-collector

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/3xui-metrics-collector .
EXPOSE 2112
CMD ["./3xui-metrics-collector"] 