package servicebuscli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/common-go/log"
)

// QueueEntity structure
type QueueEntity struct {
	Name                     string
	LockDuration             time.Duration
	AutoDeleteOnIdle         time.Duration
	DefaultMessageTimeToLive time.Duration
	MaxDeliveryCount         int32
	Forward                  ForwardEntity
	ForwardDeadLetter        ForwardEntity
}

// NewQueue Creates a Queue entity
func NewQueue(name string) QueueEntity {
	result := QueueEntity{
		MaxDeliveryCount: 10,
	}

	result.Name = name
	result.Forward.In = ForwardToQueue

	return result
}

// MapMessageForwardFlag Maps a forward flag string into it's sub components
func (s *QueueEntity) MapMessageForwardFlag(value string) {
	if value != "" {
		forwardMapped := strings.Split(value, ":")
		if len(forwardMapped) == 1 {
			s.Forward.To = forwardMapped[0]
		} else if len(forwardMapped) == 2 {
			s.Forward.To = forwardMapped[1]
			switch strings.ToLower(forwardMapped[0]) {
			case "topic":
				s.Forward.In = ForwardToTopic
			case "queue":
				s.Forward.In = ForwardToQueue
			}
		}
	}
}

// MapDeadLetterForwardFlag Maps a forward dead letter flag string into it's sub components
func (s *QueueEntity) MapDeadLetterForwardFlag(value string) {
	if value != "" {
		forwardMapped := strings.Split(value, ":")
		if len(forwardMapped) == 1 {
			s.ForwardDeadLetter.To = forwardMapped[0]
		} else if len(forwardMapped) == 2 {
			s.ForwardDeadLetter.To = forwardMapped[1]
			switch strings.ToLower(forwardMapped[0]) {
			case "topic":
				s.ForwardDeadLetter.In = ForwardToTopic
			case "queue":
				s.ForwardDeadLetter.In = ForwardToQueue
			}
		}
	}
}

// GetQueueManager creates a Service Bus Queue manager
func (s *ServiceBusCli) GetQueueManager() *servicebus.QueueManager {
	logger.Trace("Creating a service bus queue manager for service bus " + s.Namespace.Name)
	if s.Namespace == nil {
		_, err := s.GetNamespace()
		if err != nil {
			return nil
		}
	}

	s.QueueManager = s.Namespace.NewQueueManager()
	return s.QueueManager
}

// GetQueue Gets a Queue object from the Service Bus Namespace
func (s *ServiceBusCli) GetQueue(queueName string) *servicebus.Queue {
	logger.Trace("Getting queue " + queueName + " from service bus " + s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if queueName == "" {
		return nil
	}
	if s.QueueManager == nil {
		s.GetQueueManager()
	}

	qe, err := s.QueueManager.Get(ctx, queueName)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	queue, err := s.Namespace.NewQueue(qe.Name)

	return queue
}

// ListQueues Lists all the Queues in a Service Bus
func (s *ServiceBusCli) ListQueues() ([]*servicebus.QueueEntity, error) {
	logger.LogHighlight("Getting all queues from %v service bus ", log.Info, s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	qm := s.GetQueueManager()
	if qm == nil {
		commonError := errors.New("There was an error getting the queue manager, check your internet connection")
		logger.LogHighlight("There was an error getting the %v, check your internet connection", log.Error, "queue manager")
		return nil, commonError
	}

	return qm.List(ctx)
}

// CreateQueue Creates a queue in the service bus namespace
func (s *ServiceBusCli) CreateQueue(queue QueueEntity) error {
	var commonError error
	opts := make([]servicebus.QueueManagementOption, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if queue.Name == "" {
		commonError = errors.New("Queue name cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}
	logger.LogHighlight("Creating queue %v in service bus %v", log.Info, queue.Name, s.Namespace.Name)

	qm := s.GetQueueManager()

	// Checking if the queue already exists in the namespace
	existingSubscription, err := qm.Get(ctx, queue.Name)
	if existingSubscription != nil {
		commonError = errors.New("Subscription " + queue.Name + " already exists  in service bus" + s.Namespace.Name)
		logger.LogHighlight("Subscription %v already exists in service bus %v", log.Error, queue.Name, s.Namespace.Name)
		return commonError
	}

	// Generating subscription options
	if queue.LockDuration.Milliseconds() > 0 {
		opts = append(opts, servicebus.QueueEntityWithLockDuration(&queue.LockDuration))
	}
	if queue.DefaultMessageTimeToLive.Microseconds() > 0 {
		opts = append(opts, servicebus.QueueEntityWithMessageTimeToLive(&queue.DefaultMessageTimeToLive))
	}
	if queue.AutoDeleteOnIdle.Microseconds() > 0 {
		opts = append(opts, servicebus.QueueEntityWithAutoDeleteOnIdle(&queue.AutoDeleteOnIdle))
	}
	if queue.MaxDeliveryCount > 0 && queue.MaxDeliveryCount != 10 {
		opts = append(opts, servicebus.QueueEntityWithMaxDeliveryCount(int32(queue.MaxDeliveryCount)))
	}

	// Generating the forward rule, checking if the target exists or not
	if queue.Forward.To != "" {
		switch queue.Forward.In {
		case ForwardToTopic:
			tm := s.GetTopicManager()
			target, err := tm.Get(ctx, queue.Forward.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding topic %v in service bus %v", log.Error, queue.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.QueueEntityWithAutoForward(target))
		case ForwardToQueue:
			qm := s.GetQueueManager()
			target, err := qm.Get(ctx, queue.Forward.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding queue %v in service bus %v", log.Error, queue.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.QueueEntityWithForwardDeadLetteredMessagesTo(target))
		}
	}

	// Generating the Dead Letter forwarding rule, checking if the target exist or not
	if queue.ForwardDeadLetter.To != "" {
		switch queue.ForwardDeadLetter.In {
		case ForwardToTopic:
			tm := s.GetTopicManager()
			target, err := tm.Get(ctx, queue.ForwardDeadLetter.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding topic %v in service bus %v", log.Error, queue.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.QueueEntityWithAutoForward(target))
		case ForwardToQueue:
			qm := s.GetQueueManager()
			target, err := qm.Get(ctx, queue.ForwardDeadLetter.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding queue %v in service bus %v", log.Error, queue.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.QueueEntityWithAutoForward(target))
		}
	}

	_, err = qm.Put(ctx, queue.Name, opts...)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.LogHighlight("Queue %v was created successfully in service bus %v", log.Info, queue.Name, s.Namespace.Name)
	return nil
}

// DeleteQueue Deletes a queue in the service bus namespace
func (s *ServiceBusCli) DeleteQueue(queueName string) error {
	var commonError error
	if queueName == "" {
		commonError = errors.New("Queue cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	logger.LogHighlight("Removing queue %v in service bus %v", log.Info, queueName, s.Namespace.Name)
	qm := s.GetQueueManager()

	err := qm.Delete(ctx, queueName)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.LogHighlight("Queue %v was removed successfully from service bus %v", log.Info, queueName, s.Namespace.Name)
	return nil
}

// SendQueueMessage Sends a Service Bus Message to a Queue
func (s *ServiceBusCli) SendQueueMessage(queueName string, message map[string]interface{}, label string, userParameters map[string]interface{}) error {
	var commonError error
	logger.LogHighlight("Sending a service bus queue message to %v queue in service bus %v", log.Info, queueName, s.Namespace.Name)
	if queueName == "" {
		commonError = errors.New("Queue cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := s.GetQueue(queueName)
	if queue == nil {
		commonError = errors.New("Could not find queue " + queueName + " in service bus " + s.Namespace.Name)
		logger.LogHighlight("Could not find queue %v in service bus %v", log.Info, queueName, s.Namespace.Name)
		return commonError
	}

	messageData, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	sbMessage := servicebus.Message{
		Data:           messageData,
		UserProperties: userParameters,
	}

	if label != "" {
		sbMessage.Label = label
	}

	err = queue.Send(ctx, &sbMessage)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.LogHighlight("Service bus queue message was sent successfully to %v queue in service bus %v", log.Info, queueName, s.Namespace.Name)
	logger.Info("Message:")
	logger.Info(string(messageData))
	return nil
}

// SubscribeToQueue Subscribes to a queue and listen to the messages
func (s *ServiceBusCli) SubscribeToQueue(queueName string) error {
	var commonError error

	var concurrentHandler servicebus.HandlerFunc = func(ctx context.Context, msg *servicebus.Message) error {
		logger.LogHighlight("%v Received message %v on queue %v with label %v", log.Info, msg.SystemProperties.EnqueuedTime.String(), msg.ID, queueName, msg.Label)
		logger.Info("User Properties:")
		jsonString, _ := json.MarshalIndent(msg.UserProperties, "", "  ")
		fmt.Println(string(jsonString))
		logger.Info("Message Body:")
		fmt.Println(string(msg.Data))

		if !s.Peek {
			return msg.Complete(ctx)
		}
		return nil
	}

	logger.LogHighlight("Subscribing to queue %v in service bus %v", log.Info, queueName, s.Namespace.Name)
	if queueName == "" {
		commonError = errors.New("Queue " + queueName + " cannot be null")
		logger.LogHighlight("Queue %v cannot be null", log.Error, queueName)
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := s.GetQueue(queueName)
	if queue == nil {
		commonError = errors.New("Could not find queue " + queueName + " in service bus" + s.Namespace.Name)
		logger.LogHighlight("Could not find queue %v in service bus %v", log.Error, queueName, s.Namespace.Name)
		return commonError
	}
	s.ActiveQueue = queue

	logger.LogHighlight("Starting to receive messages queue %v for service bus %v", log.Info, queueName, s.Namespace.Name)
	receiver, err := queue.NewReceiver(ctx)

	if err != nil {
		commonError := errors.New("Could not create channel for queue " + queueName + " in " + s.Namespace.Name + " bus, subscription was not found")
		logger.LogHighlight("Could not create channel for queue %v for service bus %v, subscription was not found", log.Error, queueName, s.Namespace.Name)
		return commonError
	}

	listenerHandler := receiver.Listen(ctx, concurrentHandler)
	s.ActiveQueueListenerHandle = listenerHandler
	defer listenerHandler.Close(ctx)

	if <-s.CloseQueueListener {
		s.CloseQueueSubscription()
	}
	return nil
}

// CloseQueueSubscription closes the subscription to a queue
func (s *ServiceBusCli) CloseQueueSubscription() error {
	logger.LogHighlight("Closing the subscription for %v queue in service bus %v", log.Info, s.ActiveQueue.Name, s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	s.ActiveQueueListenerHandle.Close(ctx)
	s.ActiveQueue = nil
	s.ActiveQueueListenerHandle = nil
	s.CloseQueueListener <- false
	return nil
}
