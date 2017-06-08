#!/usr/bin/env bash

docker build -t webinar:debug .
docker run --rm --name=webinar-debug -p 8000:8000 -p 2345:40000 --security-opt=seccomp:unconfined webinar:debug
