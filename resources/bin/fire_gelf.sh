#!/bin/bash

CNT_IP=$(docker inspect $(docker ps -q -f name=qwatch-static) |jq -r '.[] | .NetworkSettings.Networks.bridge.IPAddress')
docker run -t --rm --name hello --log-driver gelf --log-opt gelf-address=udp://${CNT_IP}:12201 --log-opt gelf-compression-type=none debian:latest echo $@
