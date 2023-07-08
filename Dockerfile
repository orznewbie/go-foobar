FROM golang:1.17 AS builder

RUN mkdir /go-foobar
COPY . /go-foobar
WORKDIR /go-foobar/cmd/foobar
RUN GOPROXY=https://goproxy.cn go build -o foobar .

FROM busybox

COPY --from=builder /go-foobar/cmd/foobar/foobar /
CMD ["./foobar"]
