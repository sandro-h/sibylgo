VERSION=`cat version.txt`

deps:
	go get -v -t -d ./...

build: build-linux build-win

build-linux:
	go build -v -ldflags="-X 'main.buildVersion=$(VERSION)' -X 'main.buildNumber=${BUILD_NUM}' -X 'main.buildRevision=${GIT_SHA}' ${EXTRA_LDFLAGS}"

build-win:
	GOOS=windows GOARCH=amd64 go build -v -ldflags="-X 'main.buildVersion=$(VERSION)' -X 'main.buildNumber=${BUILD_NUM}' -X 'main.buildRevision=${GIT_SHA}' ${EXTRA_LDFLAGS}"

build-optimized: EXTRA_LDFLAGS=-s -w
build-optimized: build

upx:
	wget https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz
	tar xf upx-3.96-amd64_linux.tar.xz
	mv upx-3.96-amd64_linux/ upx

compress-binaries: upx
	chmod +x sibylgo
	upx/upx -q --brute sibylgo
	upx/upx -q --brute sibylgo.exe

test:
	go test -v -coverprofile="coverage.out" ./...

lint:
	golint -set_exit_status ./...

staticcheck:
	staticcheck ./...

prepare-package:
	rm -rf dist_pkg
	mkdir -p dist_pkg/sibylgo
	cp sibylgo sibylgo.exe sibylcal.html dist_pkg/sibylgo/
	cp vscode_ext/sibyl.vsix dist_pkg/sibylgo/
	chmod +x dist_pkg/sibylgo/sibylgo

	mkdir -p dist_pkg/sibylgo/outlook_cli
	cp outlook_cli/bin/Release/* dist_pkg/sibylgo/outlook_cli/

zip-package:
	rm -f sibylgo.zip
	cd dist_pkg && zip -r ../sibylgo sibylgo/*

print-version:
	echo "::set-output name=version::$$(dist_pkg/sibylgo/sibylgo --version)"

package: prepare-package zip-package print-version

release:
	git tag v$(VERSION)
	git push origin v$(VERSION)
