FROM golang:1.8.3
ADD . /go/src/github.com/dlsniper/webinar
RUN go build -gcflags="-N -l" -o /webinar github.com/dlsniper/webinar

WORKDIR /
ADD dlv /
ADD ui /ui

RUN chmod +x /dlv /webinar

EXPOSE 8000 40000

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/webinar"]
