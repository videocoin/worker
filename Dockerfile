FROM golang:1.11.1-alpine AS builder

RUN apk update && apk add --update build-base alpine-sdk musl-dev musl

WORKDIR /go/src/gitlab.videocoin.io/videocoin/transcode

ADD . ./


ENV GO111MODULE auto
RUN go build -o bin/transcoder --ldflags '-w -linkmode external -extldflags "-static"' cmd/main.go



FROM alpine:latest AS release

RUN apk update

COPY --from=builder /go/src/gitlab.videocoin.io/videocoin/transcode/bin/transcoder ./


ENTRYPOINT './transcoder'