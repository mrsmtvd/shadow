package aws

import (
	"github.com/kihamo/shadow/resource"
)

func (s *AwsService) GetConfigVariables() []resource.ConfigVariable {
	return []resource.ConfigVariable{
		resource.ConfigVariable{
			Key:   "aws-key",
			Value: "",
			Usage: "AWS access key ID",
		},
		resource.ConfigVariable{
			Key:   "aws-secret",
			Value: "",
			Usage: "AWS secret access key",
		},
		resource.ConfigVariable{
			Key:   "aws-region",
			Value: "us-east-1",
			Usage: "AWS region",
		},
	}
}
