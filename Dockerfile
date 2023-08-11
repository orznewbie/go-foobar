FROM alpine

WORKDIR /foobar

COPY ./cmd/foobar/foobar .
COPY ./etc/foobar.json .

RUN chmod +x foobar

CMD ["./foobar", "-c", "foobar.json"]
