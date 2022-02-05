VERSION=`cat version.txt`

deps:
	go get -v -t -d ./...

.PHONY: install-sys-packages
install-sys-packages:
	sudo apt update && \
	sudo apt install gcc libc6-dev \
	libx11-dev xorg-dev libxtst-dev libpng++-dev \
	xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev \
	libgl1-mesa-dev \
	gcc-mingw-w64-x86-64 libz-mingw-w64-dev

build: build-linux build-win

build-linux:
	go build -v -ldflags="-X 'main.buildVersion=$(VERSION)' -X 'main.buildNumber=${BUILD_NUM}' -X 'main.buildRevision=${GIT_SHA}' ${EXTRA_LDFLAGS}"

build-win:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -v -ldflags="-X 'main.buildVersion=$(VERSION)' -X 'main.buildNumber=${BUILD_NUM}' -X 'main.buildRevision=${GIT_SHA}' ${EXTRA_LDFLAGS}"

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

/tmp/zlib/mingw64/bin/zlib1.dll:
	sudo apt update && sudo apt install zstd
	mkdir /tmp/zlib
	wget https://mirror.msys2.org/mingw/mingw64/mingw-w64-x86_64-zlib-1.2.11-9-any.pkg.tar.zst -O /tmp/zlib.tar.zst
	tar --use-compress-program=unzstd -xvf /tmp/zlib.tar.zst -C /tmp/zlib

prepare-package: /tmp/zlib/mingw64/bin/zlib1.dll
	rm -rf dist_pkg
	mkdir -p dist_pkg/sibylgo
	cp sibylgo sibylgo.exe sibylcal.html dist_pkg/sibylgo/
	cp /tmp/zlib/mingw64/bin/zlib1.dll dist_pkg/sibylgo/
	cp vscode_ext/sibyl.vsix dist_pkg/sibylgo/
	chmod +x dist_pkg/sibylgo/sibylgo

	mkdir -p dist_pkg/sibylgo/outlook_cli
	cp outlook_cli/bin/Release/* dist_pkg/sibylgo/outlook_cli/

zip-package:
	rm -f sibylgo.zip
	cd dist_pkg && zip -r ../sibylgo sibylgo/*

print-version:
	echo "::set-output name=version::$(VERSION).${BUILD_NUM}"

package: prepare-package zip-package print-version

release:
	git tag v$(VERSION)
	git push origin v$(VERSION)
