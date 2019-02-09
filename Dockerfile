FROM ubuntu:latest AS release

LABEL maintainer="Videocoin" description="transcoding client streams"

RUN apt update && apt upgrade -y
RUN apt install ffmpeg build-essential ca-certificates -y

WORKDIR /go/src/github.com/VideoCoin/transcoder

ADD keys keys
ADD release/transcoder-linux-amd64 ./

EXPOSE 50051 50052 50053 50054 50055

ENTRYPOINT [ "./transcoder-linux-amd64" ]
