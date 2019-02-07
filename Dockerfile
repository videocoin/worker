FROM ubuntu:latest AS release

LABEL maintainer="Videocoin" description="transcoding client streams"

RUN apt update  && apt upgrade -y
RUN apt install ffmpeg imagemagick curl git -y

WORKDIR /go/src/github.com/VideoCoin/transcoder

ADD keys ./
ADD release/transcoder-linux-amd64 ./

EXPOSE 50051 50052 50053 50054 50055

ENTRYPOINT [ "./transcoder-linux-amd64" ]
