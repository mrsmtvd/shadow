package aws

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/resource"
)

type AwsService struct {
	application *shadow.Application

	Aws    *resource.Aws
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

	resourceLogger, err := a.GetResource("logger")
	if err != nil {
		return err
	}
	s.logger = resourceLogger.(*resource.Logger).Get(s.GetName())

	resourceAws, err := a.GetResource("aws")
	if err != nil {
		return err
	}
	s.Aws = resourceAws.(*resource.Aws)

	s.SNS = s.Aws.GetSNS()

	return nil
}

func (s *AwsService) Run() error {
	if s.application.HasResource("tasks") {
		tasks, _ := s.application.GetResource("tasks")
		tasks.(*resource.Dispatcher).AddNamedTask("aws.updater", s.getStatsJob)
	}

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
