GO_RICE_IMPORT := \
 -i github.com/kihamo/shadow/service/aws \
 -i github.com/kihamo/shadow/service/frontend \
 -i github.com/kihamo/shadow/service/slack \
 -i github.com/kihamo/shadow/service/system

all: clean build build-post

clean:
	rice clean $(GO_RICE_IMPORT)
	go clean

build:
	go build -a -v -ldflags '-w'

build-post:
	go get -t -v github.com/GeertJohan/go.rice/...
	rice -v $(GO_RICE_IMPORT) append --exec full

.PHONY: all clean build build-post
