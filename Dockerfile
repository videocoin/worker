FROM golang:latest AS builder

LABEL maintainer="Videocoin" description="transcoding client streams"


WORKDIR /go/src/github.com/videocoin/transcode

ADD ./ ./

RUN make build

FROM ubuntu:xenial AS release

COPY --from=builder /go/src/github.com/videocoin/transcode/bin/transcoder ./
COPY --from=builder /go/src/github.com/videocoin/transcode/keys ./keys

RUN apt-get update && apt upgrade -y

RUN apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common \
    ffmpeg


ENTRYPOINT [ "./transcoder" ]
