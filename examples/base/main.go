package main // import "github.com/kihamo/shadow/examples/base"

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kihamo/shadow"
	annotations "github.com/kihamo/shadow/components/annotations/instance"
	config "github.com/kihamo/shadow/components/config/instance"
	dashboard "github.com/kihamo/shadow/components/dashboard/instance"
	database "github.com/kihamo/shadow/components/database/instance"
	grpc "github.com/kihamo/shadow/components/grpc/instance"
	i18n "github.com/kihamo/shadow/components/i18n/instance"
	logger "github.com/kihamo/shadow/components/logger/instance"
	mail "github.com/kihamo/shadow/components/mail/instance"
	messengers "github.com/kihamo/shadow/components/messengers/instance"
	metrics "github.com/kihamo/shadow/components/metrics/instance"
	profiling "github.com/kihamo/shadow/components/profiling/instance"
	workers "github.com/kihamo/shadow/components/workers/instance"
)

var (
	build = "common"
)

func main() {
	application, err := shadow.NewApp(
		"Shadow base",
		"1.0",
		build,
		[]shadow.Component{
			annotations.NewComponent(),
			config.NewComponent(),
			dashboard.NewComponent(),
			database.NewComponent(),
			grpc.NewComponent(),
			i18n.NewComponent(),
			logger.NewComponent(),
			mail.NewComponent(),
			messengers.NewComponent(),
			metrics.NewComponent(),
			profiling.NewComponent(),
			workers.NewComponent(),
		},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = application.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
