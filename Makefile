VERSION = $(shell git describe --tags)

default: build

build: bindata.go
	go build -ldflags "-X main.version=$(VERSION)" -mod=readonly .

build_docker:
	docker build -t registry.local/cloudkeys-go-build -f Dockerfile.build .
	docker run --rm -i \
		-e "CGO_ENABLED=0" \
		-v "$(CURDIR):$(CURDIR)" \
		-w "$(CURDIR)" \
		registry.local/cloudkeys-go-build \
		make build

pre-commit: bindata.go

container: bindata.go
	docker build .

gen_css:
	lessc --verbose -x less/*.less assets/style.css

gen_js:
	coffee --compile -o assets coffee/*.coffee

bindata.go: gen_css gen_js
	go generate

publish:
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	bash golang.sh

prepare-gae-deploy:
	rm -rf Dockerfile vendor

.PHONY: bindata.go
