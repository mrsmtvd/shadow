package main

import (
	"log"

	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
	"github.com/kihamo/shadow/service/api"
	"github.com/kihamo/shadow/service/aws"
	"github.com/kihamo/shadow/service/frontend"
	"github.com/kihamo/shadow/service/slack"
	"github.com/kihamo/shadow/service/system"
	"github.com/kihamo/shadow/service/tasks"
)

func main() {
	application, err := shadow.NewApplication(
		[]shadow.Resource{
			new(resource.Config),
			new(resource.Logger),
			new(resource.Template),
		},
		[]shadow.Service{
			new(system.SystemService),
			new(tasks.TasksService),
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
