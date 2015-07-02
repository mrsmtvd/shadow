CURRENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

gen: thrift precommit

thrift:
	thrift -v -r --gen go:thrift_import=github.com/apache/thrift/lib/go/thrift,package_prefix=github.com/kihamo/shadow/service/api -o $(CURRENT_DIR)service/api $(CURRENT_DIR)service/api/service.thrift

precommit:
	goimports -w $(CURRENT_DIR)

.PHONY: gen thrift precommit