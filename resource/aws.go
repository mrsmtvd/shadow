package resource

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/kihamo/shadow"
)

type Aws struct {
	application *shadow.Application
	config      *Config
	logger      *logrus.Entry
	services    map[string]interface{}
}

type AwsArnParse struct {
	Arn          string
	Partition    string
	Service      string
	Region       string
	Account      string
	Resource     string
	ResourceType string
}

func (r *Aws) GetName() string {
	return "aws"
}

func (r *Aws) GetConfigVariables() []ConfigVariable {
	return []ConfigVariable{
		ConfigVariable{
			Key:   "aws.key",
			Value: "",
			Usage: "AWS access key ID",
		},
		ConfigVariable{
			Key:   "aws.secret",
			Value: "",
			Usage: "AWS secret access key",
		},
		ConfigVariable{
			Key:   "aws.region",
			Value: "us-east-1",
			Usage: "AWS region",
		},
	}
}

func (r *Aws) Init(a *shadow.Application) error {
	r.application = a
	r.services = map[string]interface{}{}

	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}

	r.config = resourceConfig.(*Config)

	return nil
}

func (r *Aws) Run() error {
	resourceLogger, err := r.application.GetResource("logger")
	if err != nil {
		return err
	}
	logger := resourceLogger.(*Logger).Get(r.GetName())

	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(r.config.GetString("aws.key"), r.config.GetString("aws.secret"), ""),
		Region:      aws.String(r.config.GetString("aws.region")),
	}

	if r.config.GetBool("debug") {
		awsConfig.LogLevel = aws.LogLevel(aws.LogDebug)
	}

	aws.DefaultConfig = aws.DefaultConfig.Merge(awsConfig)

	fields := logrus.Fields{
		"region": aws.DefaultConfig.Region,
	}

	credentials, err := aws.DefaultConfig.Credentials.Get()
	if err == nil {
		fields["key"] = credentials.AccessKeyID
		fields["secret"] = credentials.SecretAccessKey
	}
	logger.WithFields(fields).Info("Connect AWS")

	return nil
}

func (r *Aws) GetSNS() *sns.SNS {
	if _, ok := r.services["sns"]; !ok {
		r.services["sns"] = sns.New(aws.DefaultConfig)
	}

	return r.services["sns"].(*sns.SNS)
}

func (r *Aws) ParseArn(arn string) *AwsArnParse {
	// http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html#genref-arns

	parts := strings.Split(arn, ":")
	result := AwsArnParse{
		Arn:       parts[0],
		Partition: parts[1],
		Service:   parts[2],
		Region:    parts[3],
		Account:   parts[4],
	}

	if len(parts) > 6 {
		result.Resource = parts[5]
		result.ResourceType = parts[6]
	} else {
		path := strings.Split(parts[5], "/")

		result.Resource = path[0]
		result.ResourceType = strings.Join(path[1:], "/")
	}

	return &result
}

func (r *Aws) GetServices() map[string]interface{} {
	return r.services
}
