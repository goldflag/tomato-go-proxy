# syntax=docker/dockerfile:1

FROM golang:1.19-alpine
WORKDIR /container
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./app/main.go

EXPOSE 8000

CMD [ "/app" ]
