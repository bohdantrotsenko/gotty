OUTPUT_DIR = ./builds
GIT_COMMIT = $$(git rev-parse HEAD | cut -c1-7)
VERSION = 3.0.0-alpha.1
BUILD_OPTIONS = -ldflags "-X main.Version=$(VERSION) -X main.CommitID=$(GIT_COMMIT)"

gotty: main.go server/*.go webtty/*.go backend/*.go Makefile
	go build ${BUILD_OPTIONS}

.PHONY: all
all: gotty

server/static:
	mkdir -p server/static

server/static/index.html: server/static resources/index.html
	cp resources/index.html server/static/index.html

server/static/favicon.png: server/static resources/favicon.png
	cp resources/favicon.png server/static/favicon.png

server/static/js: server/static
	mkdir -p server/static/js


server/static/js/gotty-bundle.js: server/static/js js/dist/gotty-bundle.js
	cp js/dist/gotty-bundle.js server/static/js/gotty-bundle.js

server/static/js/gotty-bundle.js.LICENSE.txt: js/dist/gotty-bundle.js.LICENSE.txt
	cp $< $@
server/static/js/gotty-bundle.js.map: js/dist/gotty-bundle.js.map
	cp $< $@

server/static/css: server/static
	mkdir -p server/static/css

server/static/css/index.css: server/static/css resources/index.css
	cp resources/index.css server/static/css/index.css

server/static/css/xterm_customize.css: server/static/css resources/xterm_customize.css
	cp resources/xterm_customize.css server/static/css/xterm_customize.css

server/static/css/xterm.css: server/static/css js/node_modules/@xterm/xterm/css/xterm.css
	cp js/node_modules/@xterm/xterm/css/xterm.css server/static/css/xterm.css

js/node_modules/@xterm/xterm/css/xterm.css:
	cd js && \
		npm ci

js/dist/gotty-bundle.js: js/src/* js/node_modules/webpack
	cd js && \
		npx webpack

js/dist/gotty-bundle.js.LICENSE.txt: js/dist/gotty-bundle.js
js/dist/gotty-bundle.js.map: js/dist/gotty-bundle.js

js/node_modules/webpack:
	cd js && \
		npm install

tools:
	# TODO convert to `go install`
	go get github.com/mitchellh/gox           # for crosscompiling
	go get github.com/tcnksm/ghr              # for making gihub releases

.PHONY: test
test:
	if [ $$(go fmt $$(go list ./...) | wc -l) -gt 0 ]; then echo "go fmt error"; exit 1; fi
	cd js && npx tsfmt -r --verify

.PHONY: tsfmt
tsfmt:
	cd js && npx tsfmt -r

cross_compile:
	GOARM=5 gox -os="darwin linux freebsd netbsd openbsd" -arch="386 amd64 arm" -osarch="!darwin/arm" -output "${OUTPUT_DIR}/pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"

targz:
	mkdir -p ${OUTPUT_DIR}/dist
	cd ${OUTPUT_DIR}/pkg/; for osarch in *; do (cd $$osarch; tar zcvf ../../dist/gotty_${VERSION}_$$osarch.tar.gz ./*); done;

shasums:
	cd ${OUTPUT_DIR}/dist; sha256sum * > ./SHA256SUMS

release:
	ghr -c ${GIT_COMMIT} --delete --prerelease -u yudai -r gotty pre-release ${OUTPUT_DIR}/dist

.PHONY: update_go_modules
update_go_modules:
	go get -u

.PHONY: clean
clean:
	rm -rf js/node_modules/*
