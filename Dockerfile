FROM golang:1.11.1-alpine AS builder

RUN apk update && apk add --update build-base alpine-sdk musl-dev musl

WORKDIR /go/src/github.com/VideoCoin/transcode

ADD . ./


ENV GO111MODULE off
RUN make build-alpine



FROM alpine:latest AS release

RUN apk update

COPY --from=builder /go/src/github.com/VideoCoin/transcode/bin/transcoder ./


ENTRYPOINT './transcoder'