CURRENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

ifneq "$(MAKEFLAGS)" "DEBUG"
	DEBUG := false
endif

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
	cd service/api && go-bindata-assetfs -debug=$(DEBUG) -pkg="api" templates/... public/...
	cd service/aws && go-bindata-assetfs -debug=$(DEBUG) -pkg="aws" templates/...
	cd service/frontend && go-bindata-assetfs -debug=$(DEBUG) -pkg="frontend" templates/... public/...
	cd service/slack && go-bindata-assetfs -debug=$(DEBUG) -pkg="slack" templates/...
	cd service/system && go-bindata-assetfs -debug=$(DEBUG) -pkg="system" templates/...

.PHONY: gen docker format static