FROM alpine:latest

WORKDIR app

COPY broker_app .

CMD ["/app/broker_app"]