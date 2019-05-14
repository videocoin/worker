FROM golang:latest AS builder

LABEL maintainer="Videocoin" description="transcoding client streams"


WORKDIR /go/src/github.com/VideoCoin/transcode

ADD ./ ./

RUN make build

FROM ubuntu:latest AS release

COPY --from=builder /go/src/github.com/VideoCoin/transcode/bin/transcoder ./
COPY --from=builder /go/src/github.com/VideoCoin/transcode/keys ./keys

RUN apt update && apt upgrade -y
RUN apt install ca-certificates ffmpeg -y


ENTRYPOINT [ "./transcoder" ]
