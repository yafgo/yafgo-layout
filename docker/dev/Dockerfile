FROM golang:1.20-alpine

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
ENV TZ Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --update --no-cache tzdata curl jq

WORKDIR /usr/src/app
