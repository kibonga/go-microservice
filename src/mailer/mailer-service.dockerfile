FROM alpine:latest

WORKDIR /app
COPY mailer_app .
WORKDIR /app/templates
COPY templates .
WORKDIR /app

CMD ["/app/mailer_app"]

