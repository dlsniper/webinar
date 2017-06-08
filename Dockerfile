# build stage
FROM golang:1.8.3 AS build-env
ADD . /go/src/github.com/dlsniper/webinar
RUN go build -o /webinar github.com/dlsniper/webinar

# final stage
FROM ubuntu:16.04
WORKDIR /
COPY --from=build-env /webinar /

EXPOSE 8000

CMD ["/webinar"]
