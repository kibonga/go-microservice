FROM alpine:latest

WORKDIR /app

COPY logger_app .

CMD ["/app/logger_app"]