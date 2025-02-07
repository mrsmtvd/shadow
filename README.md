Shadow framework
================

[![Build Status](https://travis-ci.org/mrsmtvd/shadow.svg)](https://travis-ci.org/mrsmtvd/shadow)
[![Coverage Status](https://coveralls.io/repos/mrsmtvd/shadow/badge.svg?branch=master&service=github)](https://coveralls.io/github/mrsmtvd/shadow?branch=master)
[![GoDoc](https://godoc.org/github.com/mrsmtvd/shadow?status.svg)](https://godoc.org/github.com/mrsmtvd/shadow)

Development
------------------
```shell
brew install protobuf
brew install bower
brew install grpc
go get -u github.com/golang/protobuf/protoc-gen-go

npm install
bower install

NODE_ENV=development gulp build
NODE_ENV=development gulp watch
```

Container build
---------------
```bash
$ make build-all
```

Container upgrade
-----------------
```bash
$ docker pull mrsmtvd/shadow-full
$ docker stop shadow
$ docker rm shadow
$ docker run -d --name shadow -p 8001:8001 -p 8080:8080 mrsmtvd/shadow-full -debug=true
```

Docker restart on MacOS
-----------------------
```bash
$ boot2docker stop
$ boot2docker start
$ boot2docker ssh 'sudo /etc/init.d/docker restart'
```

Debug mode
----------
```bash
$ make DEBUG=true
```