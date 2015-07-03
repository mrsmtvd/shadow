all: clean build docker

clean:
	go clean

build:
	go build

docker:
	docker build -t shadow-full .

run:
	docker run --rm --add-host=docker:192.168.59.3 -p 8001:8001 -p 8080:8080 --name shadow-full shadow-full

.PHONY: all clean build docker run
