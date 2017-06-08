$CONTAINER_NAME = "webinar"
$CONTAINER_TAG = "dev"

docker build -f ${pwd}"/Dockerfile" `
    -t ${CONTAINER_NAME}:${CONTAINER_TAG} `
    ${pwd}

# docker run --rm --name=webinar -p 127.0.0.1:8000:8000 webinar:dev