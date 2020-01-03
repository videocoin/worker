#!/bin/bash
set -e

if [ "$1" = 'transcoderd' ]
then
    mkdir -p /env
    transinit
    source /env/init.env

    exec transcoder mine
else
    exec "$@"
fi
