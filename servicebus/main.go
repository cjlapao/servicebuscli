package servicebus

import (
	"fmt"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/common-go/log"
)

var logger = log.Get()

// ServiceBusCli Entity
type ServiceBusCli struct {
	ConnectionString          string
	Namespace                 *servicebus.Namespace
	TopicManager              *servicebus.TopicManager
	QueueManager              *servicebus.QueueManager
	ActiveTopic               *servicebus.Topic
	ActiveSubscription        *servicebus.Subscription
	ActiveQueue               *servicebus.Queue
	ActiveQueueListenerHandle *servicebus.ListenerHandle
	ActiveTopicListenerHandle *servicebus.ListenerHandle
	Peek                      bool
	UseWiretap                bool
	DeleteWiretap             bool
	CloseTopicListener        chan bool
	CloseQueueListener        chan bool
}

// NewCli creates a new ServiceBusCli
func NewCli(connectionString string) *ServiceBusCli {
	cli := ServiceBusCli{
		Peek:             false,
		UseWiretap:       false,
		DeleteWiretap:    false,
		ConnectionString: connectionString,
	}

	cli.CloseTopicListener = make(chan bool, 1)
	cli.CloseQueueListener = make(chan bool, 1)
	cli.GetNamespace()

	return &cli
}

// GetNamespace gets a new Service Bus connection namespace
func (s *ServiceBusCli) GetNamespace() (*servicebus.Namespace, error) {
	logger.Trace("Creating a service bus namespace")

	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(s.ConnectionString))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	s.Namespace = ns

	return s.Namespace, nil
}
