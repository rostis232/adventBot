FROM alpine:latest

RUN mkdir /app

COPY .  /app

WORKDIR /app

RUN chmod +x ./main

CMD ["/app/main"]