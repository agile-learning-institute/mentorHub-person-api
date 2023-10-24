# Stage 1: Build the Go binary
FROM golang:1.17 AS builder

WORKDIR /app
ADD . .
RUN GOOS=linux GOARCH=amd64 go build -o institute-person-api src/main.go

# Stage 2: Create the production image
FROM amd64/ubuntu:latest

COPY --from=builder /app/institute-person-api /institute-person-api/institute-person-api

EXPOSE 8081
CMD ["/institute-person-api/institute-person-api"]
