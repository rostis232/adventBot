#base go image

FROM golang:1.22-alpine as builder

RUN apk add --no-cache \
    gcc \
    musl-dev \
    build-base

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=1 go build -o bot ./cmd/adventbot

RUN chmod +x ./bot

#build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app /app

WORKDIR /app

CMD ["/app/bot"]