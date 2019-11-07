FROM golang:alpine as builder

COPY . /go/src/github.com/Luzifer/cloudkeys-go
WORKDIR /go/src/github.com/Luzifer/cloudkeys-go

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly

FROM alpine:latest

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/cloudkeys-go /usr/local/bin/cloudkeys-go

EXPOSE 3000
VOLUME ["/data"]

ENTRYPOINT ["/usr/local/bin/cloudkeys-go"]
CMD ["--"]

# vim: set ft=Dockerfile:
