GO_RICE_IMPORT := \
 -i github.com/kihamo/shadow/service/aws \
 -i github.com/kihamo/shadow/service/frontend \
 -i github.com/kihamo/shadow/service/slack \
 -i github.com/kihamo/shadow/service/system

all: clean build

clean:
	go clean

build: deps
	go build -a -race
	rice $(GO_RICE_IMPORT) append --exec full

deps:
	go get -t -d -v github.com/GeertJohan/go.rice/...

.PHONY: all clean build deps
