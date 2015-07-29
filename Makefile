VERSION = $(shell git describe --tags)

default: build

build: bundle_assets
	go build .

pre-commit: bundle_assets

container: ca-certificates.pem bundle_assets
	docker run -v $(CURDIR):/src -e LDFLAGS='-X main.version $(VERSION)' centurylink/golang-builder:latest
	docker build .

gen_css:
	lessc --verbose -O2 -x less/*.less assets/style.css

gen_js:
	coffee --compile -o assets coffee/*.coffee

bundle_assets: gen_css gen_js
	go-bindata assets templates

ca-certificates.pem:
		curl -s https://pki.google.com/roots.pem | grep -v "^#" | grep -v "^$$" > $@
		shasum $@
