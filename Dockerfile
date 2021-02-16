FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download

RUN go build -o file-watch .

ENTRYPOINT ["/app/file-watch"]