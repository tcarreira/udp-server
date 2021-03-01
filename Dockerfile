FROM golang:1.15-alpine as builder

COPY main.go .
RUN go build -o /udp-server \
    && chmod +x /udp-server

FROM alpine:3.13
COPY --from=builder /udp-server /usr/bin/udp-server
ENTRYPOINT ["udp-server"]
EXPOSE 1337
