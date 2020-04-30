FROM golang:1.13 AS builder

LABEL maintainer="Videocoin" description="transcoding client streams"

WORKDIR /go/src/github.com/videocoin/transcode

ADD ./ ./

RUN make build

FROM jrottenberg/ffmpeg:4.1-ubuntu AS release

COPY --from=builder /go/src/github.com/videocoin/transcode/bin/transcoder /bin/transcoder
COPY --from=builder /go/src/github.com/videocoin/transcode/capacity_test.mp4 /opt/capacity_test.mp4

RUN apt-get update

RUN apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates

EXPOSE 8888

WORKDIR /

ENTRYPOINT ["transcoder"]
CMD ["mine"]
