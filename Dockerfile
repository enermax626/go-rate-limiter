FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp ./cmd/server/main.go

#FROM alpine:latest
FROM golang:latest
WORKDIR /app
COPY --from=builder /app/myapp .
COPY --from=builder /app .

ENTRYPOINT ["./myapp"]