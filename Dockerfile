FROM golang:latest AS builder

LABEL maintainer="Videocoin" description="transcoding client streams"


WORKDIR /go/src/github.com/videocoin/transcode

ADD ./ ./

RUN make build

FROM jrottenberg/ffmpeg:4.0-ubuntu AS release

COPY --from=builder /go/src/github.com/videocoin/transcode/bin/transcoder ./

RUN apt-get update && apt upgrade -y

RUN apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates

ENTRYPOINT [ "./transcoder" ]
