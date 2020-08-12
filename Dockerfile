FROM golang:1.13 AS builder

LABEL maintainer="VideoCoin" description="Transcoding Worker"

WORKDIR /go/src/github.com/videocoin/worker

ADD ./ ./

RUN make build-linux-amd

FROM jrottenberg/ffmpeg:4.3-nvidia1804 AS release

COPY --from=builder /go/src/github.com/videocoin/worker/bin/worker-linux-amd64 /bin/worker-linux-amd64
COPY --from=builder /go/src/github.com/videocoin/worker/capacity_test.mp4 /opt/capacity_test.mp4

RUN apt-get update

RUN apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates

EXPOSE 8888

WORKDIR /

ENTRYPOINT ["worker-linux-amd64"]
CMD ["start"]
