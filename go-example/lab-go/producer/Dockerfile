FROM golang:1.17-alpine AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/main .
RUN apk add --no-cache rabbitmq-c
CMD ["./main"]

