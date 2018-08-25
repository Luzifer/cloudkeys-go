FROM golang:alpine as builder

ADD . /go/src/github.com/Luzifer/cloudkeys-go
WORKDIR /go/src/github.com/Luzifer/cloudkeys-go

RUN set -ex \
 && apk add --update git \
 && go install -ldflags "-X main.version=$(git describe --tags || git rev-parse --short HEAD || echo dev)"

FROM alpine:latest

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add ca-certificates

COPY --from=builder /go/bin/cloudkeys-go /usr/local/bin/cloudkeys-go

EXPOSE 3000

VOLUME ["/data"]

ENTRYPOINT ["/usr/local/bin/cloudkeys-go"]
CMD ["--storage=local:////data", "--password-salt=changeme", "--username-salt=changeme"]

# vim: set ft=Dockerfile:
