package main // import "github.com/kihamo/shadow/examples/base"

import (
	"log"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
	"github.com/kihamo/shadow/resource/alerts"
	"github.com/kihamo/shadow/resource/metrics"
	"github.com/kihamo/shadow/resource/workers"
	"github.com/kihamo/shadow/service/frontend"
	"github.com/kihamo/shadow/service/system"
)

func main() {
	application, err := shadow.NewApplication(
		[]shadow.Resource{
			new(resource.Config),
			new(resource.Logger),
			new(resource.Template),
			new(workers.Workers),
			new(resource.Mail),
			new(metrics.Metrics),
			new(alerts.Alerts),
		},
		[]shadow.Service{
			new(system.SystemService),
			new(frontend.FrontendService),
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
