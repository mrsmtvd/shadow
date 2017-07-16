package main // import "github.com/kihamo/shadow/examples/base"

import (
	"log"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/alerts"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/grpc"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/mail"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/profiling"
	"github.com/kihamo/shadow/components/workers"
)

func main() {
	application, err := shadow.NewApp(
		"Shadow base",
		"1.0",
		"12345-base",
		[]shadow.Component{
			new(alerts.Component),
			new(config.Component),
			new(dashboard.Component),
			new(database.Component),
			new(grpc.Component),
			new(logger.Component),
			new(mail.Component),
			new(metrics.Component),
			new(profiling.Component),
			new(workers.Component),
		},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = application.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
