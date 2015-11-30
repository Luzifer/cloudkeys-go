FROM alpine

VOLUME /data
EXPOSE 3000
ENV GOPATH /go
ENTRYPOINT ["/go/bin/cloudkeys-go"]
CMD ["--storage=local:////data", "--password-salt=changeme", "--username-salt=changeme"]

ADD . /go/src/github.com/Luzifer/cloudkeys-go
WORKDIR /go/src/github.com/Luzifer/cloudkeys-go
RUN apk --update add go git ca-certificates \
 && go get github.com/tools/godep \
 && /go/bin/godep go install -ldflags "-X main.version=$(git describe --tags)" \
 && apk --purge del git go
