FROM ubuntu:12.04

RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-get update
RUN apt-get upgrade -y

RUN apt-get install -y curl git mercurial

RUN mkdir /tmp/downloads
RUN curl -o /tmp/downloads/go1.2.linux-amd64.tar.gz -L https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz
RUN mkdir -p /opt && cd /opt && tar xfz /tmp/downloads/go1.2.linux-amd64.tar.gz

ENV GOROOT /opt/go
ENV GOPATH /root/gocode
ENV PATH /opt/go/bin:/root/gocode/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

EXPOSE 8000:8000

RUN go get code.google.com/p/gcfg
RUN go get github.com/hoisie/web
RUN go get github.com/howbazaar/loggo
RUN go get github.com/garyburd/redigo/redis

CMD ["--start"]
