FROM golang:1.13 AS builder

LABEL maintainer="VideoCoin" description="Transcoding Worker"

WORKDIR /go/src/github.com/videocoin/worker

ADD ./ ./

RUN make build

FROM jrottenberg/ffmpeg:4.1-ubuntu AS release

COPY --from=builder /go/src/github.com/videocoin/worker/bin/worker /bin/worker
COPY --from=builder /go/src/github.com/videocoin/worker/capacity_test.mp4 /opt/capacity_test.mp4

RUN apt-get update

RUN apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates

EXPOSE 8888

WORKDIR /

ENTRYPOINT ["worker"]
CMD ["start"]
