FROM golang:1.10.0-alpine3.7 AS build-env

RUN apk add --no-cache libc6-compat git gcc g++
RUN go get github.com/derekparker/delve/cmd/dlv

ADD . /go/src/github.com/dlsniper/webinar
RUN go build -gcflags "all=-N -l" -o /webinar github.com/dlsniper/webinar

# final stage
FROM alpine:3.7

RUN apk add --no-cache libc6-compat

WORKDIR /
COPY --from=build-env /webinar /
COPY --from=build-env /go/bin/dlv /
ADD ui /ui

RUN chmod +x /dlv /webinar

EXPOSE 8000 40000

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/webinar"]
