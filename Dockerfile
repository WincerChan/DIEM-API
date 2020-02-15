FROM golang:alpine

RUN apk update
RUN apk add --no-cache git

WORKDIR /hitokoto

# proxy for speed up git clone and go get
# RUN git clone https://github.com/WincerChan/DIEM-API.git /hitokoto
ADD . /hitokoto

RUN CGO_ENABLED=0 go build -o /go/bin/server
COPY ./config.yaml /etc/config.yaml
ENTRYPOINT ["/go/bin/server", "/etc/config.yaml"]
