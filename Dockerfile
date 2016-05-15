FROM golang:1.6-alpine
MAINTAINER Xan Manning <xan.manning@gmail.com>

EXPOSE 80

RUN find / -perm +6000 -type f -exec chmod a-s {} \; || true

WORKDIR /usr/src/app/

COPY . /usr/src/app/

RUN go build -o deapexer deapexer.go
RUN	mkdir -p /etc/deapexer
RUN	cp deapexer /usr/local/bin/deapexer
RUN	cp config.json /etc/deapexer/config.json

CMD ["deapexer", "-n"]
