CURRENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

gen: thrift format static

thrift:
	thrift -v -r \
        --gen go:thrift_import=github.com/apache/thrift/lib/go/thrift,package_prefix=github.com/kihamo/shadow/service/api \
        -o $(CURRENT_DIR)service/api \
        $(CURRENT_DIR)service/api/service.thrift
	rm -rf $(CURRENT_DIR)service/api/gen-go/api/api-remote

build-all: gen
	docker build -t kihamo/shadow-builder:latest docker/
	docker run --rm \
        -v "$(PWD):/src" \
        -v /var/run/docker.sock:/var/run/docker.sock \
        kihamo/shadow-builder \
        kihamo
	docker push kihamo/shadow-full

format:
	goimports -w $(CURRENT_DIR)

static:
	cd service/aws && go-bindata-assetfs -pkg="aws" templates/...
	cd service/frontend && go-bindata-assetfs -pkg="frontend" templates/... public/...
	cd service/slack && go-bindata-assetfs -pkg="slack" templates/...
	cd service/system && go-bindata-assetfs -pkg="system" templates/...

.PHONY: gen thrift precommit build-all format static