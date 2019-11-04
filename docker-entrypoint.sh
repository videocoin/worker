#!/bin/bash
set -e

if [ "$1" = 'transcoder' ]; then
    mkdir -p /env
    transinit
    source /env/init.env

    exec "$@"
fi

exec "$@"
