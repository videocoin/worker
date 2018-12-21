# Transcode

Simple transcoder for videocoin. Automatic deployment on ingest of user stream

## Build

```shell
    make
```

## Publish

Required environment variables

- PROJECT_ID
- RELEASE_BUCKET # Location of binary folder in gcloud storage

```shell
    make package
    make publish
```
