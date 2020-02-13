FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /hitokoto

# proxy for speed up git clone and go get
ENV http_proxy "http://docker.for.mac.host.internal:8002"
ENV https_proxy "http://docker.for.mac.host.internal:8002"
RUN git clone https://github.com/WincerChan/DIEM-API.git /hitokoto

RUN CGO_ENABLED=0 go build -o /go/bin/server
COPY ./config.yaml /etc/config.yaml
ENTRYPOINT ["/go/bin/server", "/etc/config.yaml"]
