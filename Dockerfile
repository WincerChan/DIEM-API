FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git 

WORKDIR /hitokoto

ADD . /hitokoto

#ENV http_proxy "http://docker.for.mac.host.internal:8002"
#ENV https_proxy "http://docker.for.mac.host.internal:8002"
#RUN git clone github.com/WincerChan/DIEM-API
#RUN CGO_ENABLED=0 go build -o /go/bin/server

RUN pwd; ls -l
RUN cat /hitokoto/config.yaml
RUN CGO_ENABLED=0 go build -o /go/bin/server
COPY ./config.yaml /etc/config.yaml
ENTRYPOINT ["/go/bin/server", "/etc/config.yaml"]

#FROM scratch AS runtime
#
#
#COPY --from=builder /go/bin/server /go/bin/server
#COPY --from=builder /hitokoto/config.yaml /config.yaml
#
#ENTRYPOINT ["/go/bin/server", "config.yaml"]
