# base go image
FROM golang:latest as builder

WORKDIR /app
COPY . .

RUN go mod download

# must contain CGO_ENABLED=0 option in order to work correctly
# broker_app is the name of the .exec file being built
# . is the relative path to the .go files compared to broker-service.old.dockerfile
RUN CGO_ENABLED=0 go build -o broker_app .
RUN chmod +x /app/broker_app

# build a tiny docker image
# code will be built on the first docker image and only the .exec will be copied to the second(smaller) docker image
FROM alpine:latest

WORKDIR /app

# /app/broker_app is the path to the .exec file from the first image
# . is the path where the .exec file will be copied (note, we used WORKDIR /app, so we already positioned ourselves there)
COPY --from=builder /app/broker_app .

CMD ["/app/broker_app"]