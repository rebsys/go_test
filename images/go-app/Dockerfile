FROM golang:latest

WORKDIR /usr/src/app

COPY app/ ./
RUN go mod tidy
RUN go build

CMD ["./go_test"]
