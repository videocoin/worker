FROM ubuntu:latest AS release

LABEL maintainer="Videocoin" description="transcoding client streams"

RUN apt update && apt upgrade -y
RUN apt install ca-certificates ffmpeg -y

WORKDIR /opt/

ADD keys keys
ADD release/transcoder-linux-amd64 ./

ENTRYPOINT [ "./transcoder-linux-amd64" ]
