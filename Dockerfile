# syntax=docker/dockerfile:1

FROM golang:latest

ENV PRODUCTION=TRUE
ENV LOGGING=TRUE
ENV NODE_TLS_REJECT_UNAUTHORIZED='0'
ENV API_KEY=b1c841697e232f05fac440ba14a09b65
ENV PORT=8000

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /tomato-go-proxy

EXPOSE 8000
CMD [ "/tomato-go-proxy" ]
