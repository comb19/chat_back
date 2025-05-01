# syntax=docker/dockerfile:1

FROM golang:1.24.2-alpine3.21

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o chat_back .

EXPOSE 8080

CMD ["./chat_back"]