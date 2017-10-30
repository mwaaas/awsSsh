FROM golang:1.9.2-alpine3.6 as binaryBuilder

RUN apk update \
    && apk add git openssh\
    && go get -u github.com/golang/dep/cmd/dep
COPY . /go/src/github.com/mwaaas/awsSsh
WORKDIR /go/src/github.com/mwaaas/awsSsh
RUN dep ensure && go build -o awsSsh


FROM alpine:latest
RUN apk --no-cache add ca-certificates openssh
WORKDIR /root/
COPY --from=binaryBuilder /go/src/github.com/mwaaas/awsSsh/awsSsh /usr/local/bin

ENTRYPOINT ["awsSsh"]

