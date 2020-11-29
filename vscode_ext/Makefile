VERSION=`cat ../version.txt`

deps:
	npm install --unsafe-perm

build:
	npm version $(VERSION) --allow-same-version && \
	npm run package
