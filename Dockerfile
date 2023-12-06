FROM golang:1.20.10-alpine3.17

RUN mkdir /app
COPY .. /app
WORKDIR /app

RUN go build -o api

CMD ["./api"]
