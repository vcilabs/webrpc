ARG executable_name=webrpc

FROM golang:1.15-alpine AS builder

ARG VERSION

ENV GOPATH /tmp/.cache

COPY . .

RUN go get -d -v  && \
    go generate ./... && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -v -o /usr/local/bin/${executable_name}

FROM alpine:3.11

LABEL type="automation" \
    name="webrpc" \
    version="v1" \
    maintainer="vision critical"

COPY --from=builder /usr/local/bin/${executable_name} /usr/local/bin/

RUN apk add curl && \
    rm -rf /var/cache/apk/* && \
    chmod -R 777 /tmp

WORKDIR /tmp