package main // import "github.com/mrsmtvd/shadow/examples/demo"

import (
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/mrsmtvd/shadow"
	_ "github.com/mrsmtvd/shadow/components/annotations/instance"
	_ "github.com/mrsmtvd/shadow/components/config/instance"
	_ "github.com/mrsmtvd/shadow/components/dashboard/instance"
	_ "github.com/mrsmtvd/shadow/components/database/instance"
	_ "github.com/mrsmtvd/shadow/components/grpc/instance"
	_ "github.com/mrsmtvd/shadow/components/i18n/instance"
	_ "github.com/mrsmtvd/shadow/components/logging/instance"
	_ "github.com/mrsmtvd/shadow/components/mail/instance"
	_ "github.com/mrsmtvd/shadow/components/messengers/instance"
	_ "github.com/mrsmtvd/shadow/components/metrics/instance"
	_ "github.com/mrsmtvd/shadow/components/ota/instance"
	_ "github.com/mrsmtvd/shadow/components/profiling/instance"
	_ "github.com/mrsmtvd/shadow/components/tracing/instance"
	_ "github.com/mrsmtvd/shadow/components/workers/instance"
	_ "github.com/mrsmtvd/shadow/examples/demo/components/demo/instance"
)

var (
	build = strconv.FormatInt(time.Now().Unix(), 10)
)

func main() {
	shadow.SetName("Demo")
	shadow.SetVersion("1.0")
	shadow.SetBuild(build)

	if err := shadow.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
