VERSION=1.0.8

deps-go:
	${GOBIN}/dep ensure --vendor-only -v

build-go: build-go-linux build-go-win

build-go-linux:
	go build -v -ldflags="-X 'main.buildVersion=${VERSION}' -X 'main.buildNumber=${BUILD_NUM}' -X 'main.buildRevision=${GIT_SHA}'"

build-go-win:
	GOOS=windows GOARCH=amd64 go build -v -ldflags="-X 'main.buildVersion=${VERSION}' -X 'main.buildNumber=${BUILD_NUM}' -X 'main.buildRevision=${GIT_SHA}'"

test-go:
	go test -v -coverprofile="coverage.out" ./...

deps-vscode:
	cd vscode_ext && \
	npm install --unsafe-perm

build-vscode:
	cd vscode_ext && \
	npm version ${VERSION} --allow-same-version && \
	node node_modules/vsce/out/vsce package -o sibyl.vsix
