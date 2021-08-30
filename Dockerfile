FROM golang:1.17 as builder

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main .

FROM maven:3.8-jdk-11 as app

COPY --from=builder /app .

CMD ["/main"]