FROM alpine:latest

WORKDIR app

COPY auth_app .

CMD ["/app/auth_app"]