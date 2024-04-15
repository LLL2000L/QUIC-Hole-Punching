FROM golang:1.20.1-alpine AS build
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /root/project
COPY . ./
RUN chmod -R 755 /root/project
RUN go mod tidy
VOLUME /go/pkg/mod
VOLUME /root/.cache/go-build
RUN go build -v -o /bin/hp

FROM alpine
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add --no-cache iptables
RUN apk update && apk add iproute2 && apk add iperf3 && rm -rf /var/cache/apk/*
RUN cp -r /usr/lib/tc /lib/
ENV PATH=${PATH}:/usr/sbin/tc
WORKDIR /bin
COPY --from=build /bin/hp ./
CMD [ "./hp","s", "c","123.56.65.101:9999"]
