#
# Build image: docker build -t irisnet/dapp-server .
#
FROM golang:1.18-alpine3.15 as builder

# Set up dependencies
ENV PACKAGES make gcc git libc-dev linux-headers bash
ARG GOPROXY=http://192.168.0.60:8081/repository/go-bianjie/,http://nexus.bianjie.ai/repository/golang-group,https://goproxy.cn,direct

WORKDIR $GOPATH/src
COPY . .

# Install minimum necessary dependencies, build binary
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk add --no-cache $PACKAGES && make install


FROM alpine:3.15

COPY --from=builder /go/bin/farm /usr/local/bin/farm
CMD ["sh","-c","farm start"]
