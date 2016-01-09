Shadow framework
================

[![Build Status](https://travis-ci.org/kihamo/shadow.svg)](https://travis-ci.org/kihamo/shadow)
[![Coverage Status](https://coveralls.io/repos/kihamo/shadow/badge.svg?branch=master&service=github)](https://coveralls.io/github/kihamo/shadow?branch=master)
[![GoDoc](https://godoc.org/github.com/kihamo/shadow?status.svg)](https://godoc.org/github.com/kihamo/shadow)

Create application
------------------
```go
package main // import "github.com/kihamo/shadow/examples/shadow-full"

import (
    "log"

    "github.com/kihamo/shadow"
    "github.com/kihamo/shadow/resource"
    "github.com/kihamo/shadow/service/api"
    "github.com/kihamo/shadow/service/aws"
    "github.com/kihamo/shadow/service/frontend"
    "github.com/kihamo/shadow/service/slack"
    "github.com/kihamo/shadow/service/system"
)

func main() {
    application, err := shadow.NewApplication(
        []shadow.Resource{
            new(resource.Config),
            new(resource.Logger),
            new(resource.Template),
            new(resource.Dispatcher),
        },
        []shadow.Service{
            new(system.SystemService),
            new(api.ApiService),
            new(aws.AwsService),
            new(frontend.FrontendService),
            new(slack.SlackService),
        },
        "1.0",
        "12345-full",
    )

    if err != nil {
        log.Fatal(err.Error())
    }

    if err = application.Run(); err != nil {
        log.Fatal(err.Error())
    }
}
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