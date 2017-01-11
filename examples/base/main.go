package main // import "github.com/kihamo/shadow/examples/base"

import (
	"log"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/alerts"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/database"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/mail"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

func main() {
	application, err := shadow.NewApp(
		[]shadow.Component{
			new(config.Component),
			new(logger.Component),
			new(metrics.Component),
			new(workers.Component),
			new(database.Component),
			new(mail.Component),
			new(alerts.Component),

			new(dashboard.Component),
		},
		"Shadow base",
		"1.0",
		"12345-base",
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = application.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
