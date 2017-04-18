#!/bin/bash

for x in $(curl -s http://localhost:9200/_cat/indices |awk '/(events|metrics|logs)/{print $3}');do
	curl -sXDELETE "http://localhost:9200/$x\?pretty"
done
