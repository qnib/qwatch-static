#!/bin/bash

docker run -d --name logger -v $(pwd)/resources/bin/:/data/bin/ \
           --log-driver gelf --log-opt gelf-address=udp://127.0.0.1:12201 \
           --log-opt gelf-compression-type=none debian:latest /data/bin/logger.sh $@
