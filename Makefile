VERSION=`cat version.txt`

deps:
	go get -v -t -d ./...

build: build-linux build-win

build-linux:
	go build -v -ldflags="-X 'main.buildVersion=$(VERSION)' -X 'main.buildNumber=${BUILD_NUM}' -X 'main.buildRevision=${GIT_SHA}'"

build-win:
	GOOS=windows GOARCH=amd64 go build -v -ldflags="-X 'main.buildVersion=$(VERSION)' -X 'main.buildNumber=${BUILD_NUM}' -X 'main.buildRevision=${GIT_SHA}'"

test:
	go test -v -coverprofile="coverage.out" ./...

lint:
	golint -set_exit_status ./...

release:
	git tag v$(VERSION)
	git push origin v$(VERSION)
