$CONTAINER_NAME = "webinar"
$CONTAINER_TAG = "debug"

docker build -f ${pwd}"/Debug-Dockerfile" `
    -t ${CONTAINER_NAME}:${CONTAINER_TAG} `
    ${pwd}

# docker run --rm --name=webinar-debug -p 127.0.0.1:8000:8000 -p 127.0.0.1:2345:40000 --security-opt=seccomp:unconfined webinar:debug