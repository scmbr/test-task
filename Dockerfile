FROM golang:1.23-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/app/main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /auth-service .
COPY --from=builder /app/configs ./configs
EXPOSE 8080
CMD ["./auth-service"]