GO_RICE_IMPORT := \
 -i github.com/kihamo/shadow/service/aws \
 -i github.com/kihamo/shadow/service/frontend \
 -i github.com/kihamo/shadow/service/slack \
 -i github.com/kihamo/shadow/service/system

docker-run:
	docker run -p 8080:8080 -i -t full

build-post:
	go get -t -v github.com/GeertJohan/go.rice/...
	rice -v $(GO_RICE_IMPORT) append --exec shadow-full

.PHONY: build-post