FROM golang

RUN mkdir -p /go/src/github.com/mikolajb/pubsub-test/
WORKDIR /go/src/github.com/mikolajb/pubsub-test/

COPY . /go/src/github.com/mikolajb/pubsub-test/
# RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
