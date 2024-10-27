# Use official Golang image as build stage
FROM golang:1.23.1 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o card-service .

#RUN go build -o card-service .

# Use minimal Alpine image for final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/card-service .
RUN chmod +x /root/card-service


EXPOSE 8080

CMD ["./card-service"]

