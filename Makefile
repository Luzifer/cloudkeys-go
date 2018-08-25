VERSION = $(shell git describe --tags)

default: build

build: bindata.go
	go build -ldflags "-X main.version=$(VERSION)" .

container: bindata.go
	docker build .

build_vue:
	docker run --rm -i \
		-v "$(CURDIR):/src" \
		-w "/src" \
		node:10-alpine \
		sh -exc "npm ci && npm run build && chown -R $(shell id -u):$(shell id -g) ."

bindata.go: build_vue
	go-bindata -o bindata.go assets/...

publish:
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	bash golang.sh

prepare-gae-deploy:
	rm -rf Dockerfile vendor

.PHONY: bindata.go
