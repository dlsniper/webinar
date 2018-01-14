#!/usr/bin/env bash

proj=${GOPATH}/src/github.com/dlsniper/webinar
cd ${proj}/db

docker build -t webinar-db:db .
docker run -d --rm --name=webinar-db -p 5432:5432 webinar-db:db

cd ${proj}
docker build -t webinar-debug:dev .
docker run -d --rm --name=webinar-db -p 8000:8000 -p 40000:40000 --link webinar-db:db --security-opt="apparmor=unconfined" --cap-add=SYS_PTRACE webinar-debug:dev
