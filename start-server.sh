#!/bin/bash
#
#  Simple script that starts the server
#
#  Created: 19. Nov 2017
#

which docker >> /dev/null 2>&1
if [[ "$?" -eq 1 ]]; then
    echo "Expecting docker binaries present on the system, exiting ..."
    exit
fi

ln=$(docker images | grep 'hn0stuff/go_gdal' | wc -l)
if [[ "$ln" != 1 ]]; then
    docker pull hn0stuff/go_gdal:1.0
fi

docker run --rm -d \
           -p 8000:8000 \
           -v $(dirname $(readlink -e "$0")):/usr/src/myapp \
           -w "/usr/src/myapp/server" \
           -e GOPATH="/usr/src/myapp/server" \
           hn0stuff/go_gdal:1.0 go run src/app/app.go ../data/sample_data.sqlite ../client