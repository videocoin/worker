FROM golang:latest AS builder

WORKDIR /opt/

ADD . ./

ENV GO111MODULE off

RUN make build


FROM ubuntu:latest AS release

COPY --from=builder /opt/bin/transcoder ./

ENTRYPOINT './transcoder'