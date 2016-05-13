FROM golang:alpine

MAINTAINER Knut Ahlers <knut@ahlers.me>

ADD . /go/src/github.com/Luzifer/cloudkeys-go
WORKDIR /go/src/github.com/Luzifer/cloudkeys-go

RUN set -ex \
 && apk add --update git \
 && go install -ldflags "-X main.version=$(git describe --tags || git rev-parse --short HEAD || echo dev)" \
 && apk del --purge git

EXPOSE 3000

VOLUME ["/data"]

ENTRYPOINT ["/go/bin/cloudkeys-go"]
CMD ["--storage=local:////data", "--password-salt=changeme", "--username-salt=changeme"]
