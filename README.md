Shadow framework
================

[![Build Status](https://travis-ci.org/kihamo/shadow.svg)](https://travis-ci.org/kihamo/shadow)
[![Coverage Status](https://coveralls.io/repos/kihamo/shadow/badge.svg?branch=master&service=github)](https://coveralls.io/github/kihamo/shadow?branch=master)
[![GoDoc](https://godoc.org/github.com/kihamo/shadow?status.svg)](https://godoc.org/github.com/kihamo/shadow)

Development
------------------
```shell
brew install protobuf
brew install bower
brew install grpc

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
$ docker pull kihamo/shadow-full
$ docker stop shadow
$ docker rm shadow
$ docker run -d --name shadow -p 8001:8001 -p 8080:8080 kihamo/shadow-full -debug=true
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