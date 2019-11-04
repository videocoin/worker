FROM golang:1.12.4 AS builder

LABEL maintainer="Videocoin" description="transcoding client streams"

WORKDIR /go/src/github.com/videocoin/transcode

ADD ./ ./

RUN make build
RUN make build-transinit

FROM jrottenberg/ffmpeg:4.1-ubuntu AS release

COPY --from=builder /go/src/github.com/videocoin/transcode/bin/transinit /bin/transinit
COPY --from=builder /go/src/github.com/videocoin/transcode/bin/transcoder /bin/transcoder
COPY --from=builder /go/src/github.com/videocoin/transcode/docker-entrypoint.sh /docker-entrypoint.sh

RUN chmod a+x /docker-entrypoint.sh

RUN apt-get update

RUN apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates

WORKDIR /

RUN mkdir -p /env
RUN touch /env/init.env
RUN echo "LOGLEVEL=debug" > /env/init.env

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["transcoder"]
