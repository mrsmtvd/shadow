package main // import "github.com/kihamo/shadow/examples/full"

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
