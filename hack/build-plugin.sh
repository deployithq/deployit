#!/bin/bash
docker build -f ./images/plugins/docker/Dockerfile -t index.lstbknd.net/lastbackend/plugin .

id=$(docker create index.lstbknd.net/lastbackend/plugin true)
mkdir -p ./build/plugins/docker/rootfs
docker export "$id" | tar -x -C ./build/plugins/docker/rootfs
cp ./images/plugins/docker/config.json ./build/plugins/docker/config.json
docker rm -vf "$id"
docker plugin disable lastbackend 
docker plugin rm lastbackend
docker plugin create lastbackend ./build/plugins/docker
docker plugin enable lastbackend 