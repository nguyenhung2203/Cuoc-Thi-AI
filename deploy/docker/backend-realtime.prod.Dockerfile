FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/realtime ./cmd/realtime/main.go

FROM alpine:3.22
RUN addgroup -S app && adduser -S -G app app && apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=builder /out/realtime ./main
RUN chown -R app:app /app
USER app
EXPOSE 8081
HEALTHCHECK --interval=10s --timeout=5s --start-period=20s --retries=5 CMD wget -qO- http://127.0.0.1:8081/health || exit 1
ENTRYPOINT ["./main"]
