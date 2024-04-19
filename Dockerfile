# 可以替换成自己的镜像源, 注意最后需要一个斜杠, 如: registry.cn-shanghai.aliyuncs.com/{你的命名空间}/
ARG REGISTRY=docker.io/library/

# build
FROM ${REGISTRY}golang:1.21-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags "-s -w" -o /app/app ./cmd/server/main.go

# release
FROM ${REGISTRY}alpine:3.16

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app

# golang
COPY --from=builder /app/app /app/yafgo
COPY --from=builder /build/resource /app/resource
COPY --from=builder /build/config /app/config

ENTRYPOINT [ "/app/yafgo" ]
CMD ["-c=prod"]
