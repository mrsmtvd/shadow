CURRENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

gen: format static

docker: gen
	docker pull kihamo/go-builder
	docker run --rm \
        -v "$(PWD):/src" \
        -v /var/run/docker.sock:/var/run/docker.sock \
        kihamo/go-builder \
        kihamo
	docker push kihamo/shadow-full

format:
	goimports -w $(CURRENT_DIR)

static:
	cd service/aws && go-bindata-assetfs -pkg="aws" templates/...
	cd service/frontend && go-bindata-assetfs -pkg="frontend" templates/... public/...
	cd service/slack && go-bindata-assetfs -pkg="slack" templates/...
	cd service/system && go-bindata-assetfs -pkg="system" templates/...

.PHONY: gen docker format static