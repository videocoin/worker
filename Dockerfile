FROM golang:1.11.1-alpine AS builder

RUN apk update && apk add --update build-base alpine-sdk musl-dev musl

WORKDIR /go/src/gitlab.videocoin.io/videocoin/transcode

ADD . ./


ENV GO111MODULE auto
RUN make build-alpine



FROM alpine:latest AS release

RUN apk update

COPY --from=builder /go/src/gitlab.videocoin.io/videocoin/transcode/bin/transcoder ./


ENTRYPOINT './transcoder'