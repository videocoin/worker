FROM golang:latest AS builder

LABEL maintainer="Videocoin" description="transcoding client streams"

WORKDIR /go/src/github.com/videocoin/transcode

ADD ./ ./

ENV GO111MODULE=off

RUN make build

FROM jrottenberg/ffmpeg:4.1-ubuntu AS release

COPY --from=builder /go/src/github.com/videocoin/transcode/bin/transcoder /transcoder
COPY --from=builder /go/src/github.com/videocoin/transcode/entrypoint.sh /entrypoint.sh

RUN chmod a+x /entrypoint.sh

RUN apt-get update && apt upgrade -y

RUN apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates

ENTRYPOINT ["/entrypoint.sh"]
