#!/bin/bash

docker run -t --rm --name hello --log-driver gelf --log-opt gelf-address=udp://12201:12201 --log-opt gelf-compression-type=none debian:latest echo $@


