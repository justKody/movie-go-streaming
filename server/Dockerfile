# Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server .

# Runtime
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
	&& adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder /server .

USER appuser

EXPOSE 8080

CMD ["./server"]
