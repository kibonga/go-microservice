# base go image
FROM golang:latest as builder

WORKDIR /app

ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

#COPY go.mod go.sum ./
#COPY *.go ./
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /broker_app .

#RUN chmod +x /app/broker_app


# build a tiny docker image
# code will be built on the first docker image and only the .exec will be copied to the second(smaller) docker image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /broker_app ./

CMD ["/app/broker_app"]