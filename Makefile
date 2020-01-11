VERSION = $(shell git describe --tags)

default: build

build: bindata.go
	go build -ldflags "-X main.version=$(VERSION)" .

container: bindata.go
	docker build .

vue_dev:
	docker run --rm -i \
		-v "$(CURDIR):/src" \
		-w "/src" -u $(shell id -u) \
		-p 8080:8080 \
		node:10-alpine \
		sh -exc "apk add python && npm ci && npm run serve"

build_vue:
	docker run --rm -i \
		-v "$(CURDIR):/src" \
		-w "/src" \
		node:10-alpine \
		sh -exc "apk add python && npm ci && npm run build"

lint:
	docker run --rm -i \
		-v "$(CURDIR):/src" \
		-w "/src" \
		node:10-alpine \
		sh -exc "apk add python && npm ci && npx eslint src"

lint-fix:
	docker run --rm -i \
		-v "$(CURDIR):/src" \
		-w "/src" \
		node:10-alpine \
		sh -exc "apk add python && npm ci && npx eslint --fix src"

lint-watch:
	docker run --rm -i \
		-v "$(CURDIR):/src" \
		-w "/src" \
		node:10-alpine \
		sh -exc "apk add python && npm ci && while true; do npx eslint src || true; sleep 5; done"

.PHONY: bindata.go
bindata.go: build_vue
	go-bindata -o bindata.go dist/...

publish:
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	bash golang.sh

prepare-gae-deploy:
	rm -rf Dockerfile vendor

.PHONY: public/wasm_exec.js
public/wasm_exec.js:
	curl -sSfLo public/wasm_exec.js "https://raw.githubusercontent.com/golang/go/go1.11/misc/wasm/wasm_exec.js"

.PHONY: public/cryptocore.wasm
public/cryptocore.wasm:
	GOOS=js GOARCH=wasm go build -o public/cryptocore.wasm ./cryptocore
