docker:
	docker pull kihamo/go-builder
	docker run --rm \
        -v "$(PWD):/src" \
        -v /var/run/docker.sock:/var/run/docker.sock \
        kihamo/go-builder \
        kihamo
	docker push kihamo/shadow-full

.PHONY: docker