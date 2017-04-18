#!/bin/bash

log=23
event=42

while [ true ];do
	r=$((${RANDOM}%5000))
    echo "put random_number $(date +%s) $r script=logger.sh,hostname=$(hostname)"
    if [[ $(($r%${log})) -eq 0 ]];then
        echo "${r} is a multiple of ${log} - Let's Log!"
    fi
    if [[ $(($r%${event})) -eq 0 ]];then
		msg="${r} is a multiple of ${event}"
        ec="001.${event}"
        echo "cee{\"event_code\": \"${ec}\", \"msg\": \"${msg}\"}"
    fi
    sleep ${1:-5}
done
