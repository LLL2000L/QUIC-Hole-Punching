FROM golang:1.20.1-alpine AS build
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /root/project
COPY . ./
RUN chmod -R 755 /root/project
RUN go mod tidy
VOLUME /go/pkg/mod
VOLUME /root/.cache/go-build
RUN go build -v -o /bin/cs
RUN chmod -R 755 /bin/cs

FROM alpine
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add --no-cache iptables
RUN apk update && apk add iproute2 && apk add iperf3 && rm -rf /var/cache/apk/*
RUN cp -r /usr/lib/tc /lib/
ENV PATH=${PATH}:/usr/sbin/tc
WORKDIR /bin
COPY --from=build /bin/cs ./

CMD [ "./cs", "tcp_s", "tcp_c","192.168.100.100:9001"]


#在最开始构建 docker-compose build 的时候，由于还没有下载安装包，因此先需要执行 go mod tidy，
#之后再次构建 docker-compose build ，可以将 go mod tidy 注释掉

#RUN apk update && apk add iproute2 && apk add iperf3 && rm -rf /var/cache/apk/*
#
## 修复 iproute2 包的bug，见 https://gitlab.alpinelinux.org/alpine/aports/issues/6948
#RUN cp -r /usr/lib/tc /lib/
