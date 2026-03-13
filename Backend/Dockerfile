FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ulam-api ./cmd/api

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/ulam-api /ulam-api
EXPOSE 8080
ENTRYPOINT ["/ulam-api"]
