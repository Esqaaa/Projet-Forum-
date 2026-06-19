FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o forum .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/forum .
COPY templates/ templates/
COPY static/ static/
EXPOSE 8080
CMD ["./forum"]