FROM golang:latest AS builder

LABEL maintainer="Videocoin" description="transcoding client streams"


WORKDIR /go/src/github.com/VideoCoin/transcode

ADD ./ ./

RUN make build

FROM ubuntu:xenial AS release

COPY --from=builder /go/src/github.com/VideoCoin/transcode/bin/transcoder ./
COPY --from=builder /go/src/github.com/VideoCoin/transcode/keys ./keys

RUN apt-get update && apt upgrade -y

RUN apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common \
    ffmpeg


ENTRYPOINT [ "./transcoder" ]
