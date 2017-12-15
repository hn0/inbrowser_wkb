#!/bin/bash
#
#  Simple script that starts the server
#
#  Created: 19. Nov 2017
#

# TODO: add code for docker container
#  for first iter local installation of go will be enough

SERVER_PATH=$PWD"/server"

if [ "$GOPATH" != "$SERVER_PATH" ]; then 
    export GOPATH="$SERVER_PATH"
fi 

cd $SERVER_PATH
echo "Starting server ..."

go run src/app/app.go ../data/sample_data.sqlite

# echo "Server running"