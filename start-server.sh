#!/bin/bash
#
#  Simple script that starts the server
#
#  Created: 19. Nov 2017
#


SERVER_PATH=$PWD"/server"

if [ "$GOPATH" != "$SERVER_PATH" ]; then 
    export GOPATH="$SERVER_PATH"
fi 

cd $SERVER_PATH
echo "Starting server ..."


# usage app.go geometry_dataset static_files_folder
go run src/app/app.go ../data/sample_data.sqlite ../client

