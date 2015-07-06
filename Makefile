CURRENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

gen: thrift precommit

thrift:
	thrift -v -r \
        --gen go:thrift_import=github.com/apache/thrift/lib/go/thrift,package_prefix=github.com/kihamo/shadow/service/api \
        -o $(CURRENT_DIR)service/api \
        $(CURRENT_DIR)service/api/service.thrift
	rm -rf $(CURRENT_DIR)service/api/gen-go/api/api-remote

precommit:
	goimports -w $(CURRENT_DIR)

builder:
	docker build -t kihamo/shadow-builder:latest docker/

build-all: builder
	docker run --rm \
        -v "$(PWD):/src" \
        -v /var/run/docker.sock:/var/run/docker.sock \
        kihamo/shadow-builder \
        kihamo

.PHONY: gen thrift precommit builder build-all