FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git 

WORKDIR /hitokoto

ADD . /hitokoto

RUN go get -d -v 
RUN CGO_ENABLED=0 go build -o /go/bin/server

FROM scratch AS runtime


COPY --from=builder /go/bin/server /go/bin/server
COPY --from=builder /hitokoto/config.json /config.json

ENTRYPOINT ["/go/bin/server", "prod", "--config=config.json"]
