FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git 

WORKDIR /hitokoto

ADD . /hitokoto

RUN go get -d -v github.com/WincerChan/DIEM-API
RUN CGO_ENABLED=0 go build -o /go/bin/server

FROM scratch AS runtime


COPY --from=builder /go/bin/server /go/bin/server
COPY --from=builder /hitokoto/config.yaml /config.yaml

ENTRYPOINT ["/go/bin/server", "config.yaml"]
