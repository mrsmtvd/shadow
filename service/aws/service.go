package aws

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type AwsService struct {
	application *shadow.Application

	SNS    *sns.SNS
	logger *logrus.Entry

	mutex sync.RWMutex

	applications  []*sns.PlatformApplication
	subscriptions []*sns.Subscription
	topics        []*sns.Topic
}

func (s *AwsService) GetName() string {
	return "aws"
}

func (s *AwsService) Init(a *shadow.Application) error {
	s.application = a

	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	config := resourceConfig.(*resource.Config)

	resourceLogger, err := a.GetResource("logger")
	if err != nil {
		return err
	}
	s.logger = resourceLogger.(*resource.Logger).Get(s.GetName())

	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(config.GetString("aws-key"), config.GetString("aws-secret"), ""),
		Region:      config.GetString("aws-region"),
	}

	if config.GetBool("debug") {
		awsConfig.LogLevel = 5
	}

	s.SNS = sns.New(awsConfig)

	if a.HasResource("tasks") {
		tasks, _ := a.GetResource("tasks")
		tasks.(*resource.Dispatcher).AddTask(s.getStatsJob)
	}

	return nil
}

func (s *AwsService) Run() error {
    s.mutex.Lock()
	s.logger.Info("Connect AWS SNS")
    s.mutex.Unlock()

	return nil
}

func (s *AwsService) getStatsJob(args ...interface{}) (bool, time.Duration) {
	var stop bool

	// applications
	applications := []*sns.PlatformApplication{}
	paramsApplications := &sns.ListPlatformApplicationsInput{}
	for !stop {
		responseApps, err := s.SNS.ListPlatformApplications(paramsApplications)
		if err == nil {
			applications = append(applications, responseApps.PlatformApplications...)

			if responseApps.NextToken != nil {
				paramsApplications.NextToken = responseApps.NextToken
			} else {
				stop = true
			}
		} else {
            s.mutex.Lock()
			s.logger.Panicf(err.Error())
            s.mutex.Unlock()
			stop = true
		}
	}

	// subscriptions
	stop = false
	subscriptions := []*sns.Subscription{}
	paramsSubscriptions := &sns.ListSubscriptionsInput{}
	for !stop {
		responseSubscriptions, err := s.SNS.ListSubscriptions(paramsSubscriptions)
		if err == nil {
			subscriptions = append(subscriptions, responseSubscriptions.Subscriptions...)

			if responseSubscriptions.NextToken != nil {
				paramsSubscriptions.NextToken = responseSubscriptions.NextToken
			} else {
				stop = true
			}
		} else {
            s.mutex.Lock()
			s.logger.Panicf(err.Error())
            s.mutex.Unlock()
			stop = true
		}
	}

	// topics
	stop = false
	topics := []*sns.Topic{}
	paramsTopics := &sns.ListTopicsInput{}
	for !stop {
		responseTopics, err := s.SNS.ListTopics(paramsTopics)
		if err == nil {
			topics = append(topics, responseTopics.Topics...)

			if responseTopics.NextToken != nil {
				paramsTopics.NextToken = responseTopics.NextToken
			} else {
				stop = true
			}
		} else {
            s.mutex.Lock()
			s.logger.Panicf(err.Error())
            s.mutex.Unlock()
			stop = true
		}
	}

	s.mutex.Lock()
	s.applications = applications
	s.subscriptions = subscriptions
	s.topics = topics
	s.mutex.Unlock()

	return true, time.Hour
}
