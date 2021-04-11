package servicebuscli

import (
	"fmt"

	"github.com/cjlapao/common-go/log"

	servicebus "github.com/Azure/azure-service-bus-go"
)

// ForwardingDestination Enum
type ForwardingDestination int

// ForwardingDestination Enum definition
const (
	ForwardToTopic ForwardingDestination = iota
	ForwardToQueue
)

// ForwardEntity struct
type ForwardEntity struct {
	To string
	In ForwardingDestination
}

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

var serviceBusCli *ServiceBusCli
var logger = log.Get()

// Get creates a new ServiceBusCli
func Get(connectionString string) *ServiceBusCli {
	if serviceBusCli != nil {
		return serviceBusCli
	}

	serviceBusCli = &ServiceBusCli{
		Peek:             false,
		UseWiretap:       false,
		DeleteWiretap:    false,
		ConnectionString: connectionString,
	}

	serviceBusCli.CloseTopicListener = make(chan bool, 1)
	serviceBusCli.CloseQueueListener = make(chan bool, 1)
	serviceBusCli.GetNamespace()

	return serviceBusCli
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
