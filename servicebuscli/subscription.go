package servicebuscli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/common-go/log"
)

// SubscriptionEntity structure
type SubscriptionEntity struct {
	Name                     string
	TopicName                string
	LockDuration             time.Duration
	AutoDeleteOnIdle         time.Duration
	DefaultMessageTimeToLive time.Duration
	MaxDeliveryCount         int32
	Forward                  ForwardEntity
	ForwardDeadLetter        ForwardEntity
	Rules                    []RuleEntity
}

// RuleEntity structure
type RuleEntity struct {
	Name      string
	SQLFilter string
	SQLAction string
}

// NewSubscription Creates a new subscription entity
func NewSubscription(topicName string, name string) SubscriptionEntity {
	result := SubscriptionEntity{
		Name:             name,
		TopicName:        topicName,
		MaxDeliveryCount: 10,
	}

	result.Rules = make([]RuleEntity, 0)
	result.Forward.In = ForwardToTopic

	return result
}

// AddSQLFilter Adds a Sql filter to a specific Rule
func (s *SubscriptionEntity) AddSQLFilter(ruleName string, filter string) {
	var rule RuleEntity
	ruleFound := false

	for i := range s.Rules {
		if s.Rules[i].Name == ruleName {
			ruleFound = true
			if len(s.Rules[i].SQLFilter) > 0 {
				s.Rules[i].SQLFilter += " "
			}
			s.Rules[i].SQLFilter += filter
			break
		}
	}

	if !ruleFound {
		rule = RuleEntity{
			Name:      ruleName,
			SQLFilter: filter,
		}
		s.Rules = append(s.Rules, rule)
	}
}

// AddSQLAction Adds a Sql Action to a specific rule
func (s *SubscriptionEntity) AddSQLAction(ruleName string, action string) {
	var rule RuleEntity
	ruleFound := false

	for i := range s.Rules {
		if s.Rules[i].Name == ruleName {
			ruleFound = true
			if len(s.Rules[i].SQLAction) > 0 {
				s.Rules[i].SQLAction += " "
			}
			s.Rules[i].SQLAction += action
			break
		}
	}

	if !ruleFound {
		rule = RuleEntity{
			Name:      ruleName,
			SQLAction: action,
		}
		s.Rules = append(s.Rules, rule)
	}
}

// MapMessageForwardFlag Maps a forward flag string into it's sub components
func (s *SubscriptionEntity) MapMessageForwardFlag(value string) {
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
func (s *SubscriptionEntity) MapDeadLetterForwardFlag(value string) {
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

// MapRuleFlag Maps a rule flag string into it's sub components
func (s *SubscriptionEntity) MapRuleFlag(value string) {
	if value != "" {
		ruleMapped := strings.Split(value, ":")
		if len(ruleMapped) > 1 {
			s.AddSQLFilter(ruleMapped[0], ruleMapped[1])
			if len(ruleMapped) == 3 {
				s.AddSQLAction(ruleMapped[0], ruleMapped[2])
			}
		}
	}
}

// ListSubscriptions Lists all the topics in a service bus
func (s *ServiceBusCli) ListSubscriptions(topicName string) ([]*servicebus.SubscriptionEntity, error) {
	logger.LogHighlight("Getting all topics from %v service bus ", log.Info, s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	topic := s.GetTopic(topicName)
	if topic == nil {
		commonError := errors.New("Could not find topic " + topicName + " in " + s.Namespace.Name + " bus")
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Error, topicName, s.Namespace.Name)
		return nil, commonError
	}

	sm := topic.NewSubscriptionManager()
	return sm.List(ctx)
}

// CreateSubscription Creates a subscription to a topic in the service bus
func (s *ServiceBusCli) CreateSubscription(subscription SubscriptionEntity) error {
	var commonError error
	opts := make([]servicebus.SubscriptionManagementOption, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	logger.LogHighlight("Creating subscription %v on topic %v in service bus %v", log.Info, subscription.Name, subscription.TopicName, s.Namespace.Name)
	topic := s.GetTopic(subscription.TopicName)
	if topic == nil {
		commonError = errors.New("Could not find topic " + subscription.TopicName + " in service bus" + s.Namespace.Name)
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Error, subscription.TopicName, s.Namespace.Name)
		return commonError
	}
	sm := topic.NewSubscriptionManager()
	existingSubscription, err := sm.Get(ctx, subscription.Name)
	if existingSubscription != nil {
		commonError = errors.New("Subscription " + subscription.Name + " already exists on topic " + subscription.TopicName + " in service bus" + s.Namespace.Name)
		logger.LogHighlight("Subscription %v already exists on topic %v in service bus %v", log.Error, subscription.Name, subscription.TopicName, s.Namespace.Name)
		return commonError
	}

	// Generating subscription options
	if subscription.LockDuration.Milliseconds() > 0 {
		opts = append(opts, servicebus.SubscriptionWithLockDuration(&subscription.LockDuration))
	}
	if subscription.DefaultMessageTimeToLive.Microseconds() > 0 {
		opts = append(opts, servicebus.SubscriptionWithMessageTimeToLive(&subscription.DefaultMessageTimeToLive))
	}
	if subscription.AutoDeleteOnIdle.Microseconds() > 0 {
		opts = append(opts, servicebus.SubscriptionWithAutoDeleteOnIdle(&subscription.AutoDeleteOnIdle))
	}

	// Generating the forward rule, checking if the target exists or not
	if subscription.Forward.To != "" {
		switch subscription.Forward.In {
		case ForwardToTopic:
			tm := s.GetTopicManager()
			target, err := tm.Get(ctx, subscription.Forward.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding topic %v in service bus %v", log.Error, subscription.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.SubscriptionWithAutoForward(target))
		case ForwardToQueue:
			qm := s.GetQueueManager()
			target, err := qm.Get(ctx, subscription.Forward.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding queue %v in service bus %v", log.Error, subscription.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.SubscriptionWithAutoForward(target))
		}
	}

	// Generating the Dead Letter forwarding rule, checking if the target exist or not
	if subscription.ForwardDeadLetter.To != "" {
		switch subscription.ForwardDeadLetter.In {
		case ForwardToTopic:
			tm := s.GetTopicManager()
			target, err := tm.Get(ctx, subscription.ForwardDeadLetter.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding topic %v in service bus %v", log.Error, subscription.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.SubscriptionWithAutoForward(target))
		case ForwardToQueue:
			qm := s.GetQueueManager()
			target, err := qm.Get(ctx, subscription.ForwardDeadLetter.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding queue %v in service bus %v", log.Error, subscription.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.SubscriptionWithForwardDeadLetteredMessagesTo(target))
		}
	}

	_, err = sm.Put(ctx, subscription.Name, opts...)
	if err != nil {
		logger.Info("There was an error creating subscription")
		logger.Error(err.Error())
		return err
	}

	// Defining the filters if they exist
	for _, rule := range subscription.Rules {
		s.CreateSubscriptionRule(subscription, rule)
		if err != nil {
			return err
		}

	}

	logger.LogHighlight("Subscription %v was created successfully on topic %v in service bus %v", log.Info, subscription.Name, subscription.TopicName, s.Namespace.Name)
	return nil
}

// CreateSubscriptionRule Creates a rule to a specific subscription
func (s *ServiceBusCli) CreateSubscriptionRule(subscription SubscriptionEntity, rule RuleEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	logger.LogHighlight("Creating subscription rule %v in subscription %v on topic %v in service bus %v", log.Info, rule.Name, subscription.Name, subscription.TopicName, s.Namespace.Name)
	topic := s.GetTopic(subscription.TopicName)
	sm := topic.NewSubscriptionManager()

	var sqlFilter servicebus.SQLFilter
	var sqlAction servicebus.SQLAction
	if rule.SQLFilter != "" {
		sqlFilter.Expression = rule.SQLFilter
		if rule.SQLAction != "" {
			sqlAction.Expression = rule.SQLAction
			_, err := sm.PutRuleWithAction(ctx, subscription.Name, rule.Name, sqlFilter, sqlAction)
			if err != nil {
				logger.LogHighlight("Could not create subscription rule %v in subscription %v on topic %v in service bus %v", log.Error, rule.Name, subscription.Name, subscription.TopicName, s.Namespace.Name)
				return err
			}
		} else {
			_, err := sm.PutRule(ctx, subscription.Name, rule.Name, sqlFilter)
			if err != nil {
				logger.LogHighlight("Could not create subscription rule %v in subscription %v on topic %v in service bus %v", log.Error, rule.Name, subscription.Name, subscription.TopicName, s.Namespace.Name)
				return err
			}
		}
	}
	logger.LogHighlight("Subscription rule %v was created successfully for subscription %v on topic %v in service bus %v", log.Info, rule.Name, subscription.Name, subscription.TopicName, s.Namespace.Name)

	rules, err := sm.ListRules(ctx, subscription.Name)
	if err != nil {
		logger.LogHighlight("There was an error trying to list the rules of subscription %v on topic %v in service bus %v", log.Error, subscription.Name, subscription.TopicName, s.Namespace.Name)
		return err
	}

	for _, existingRule := range rules {
		if existingRule.Name == "$Default" {
			if len(rules) > 1 {
				sm.DeleteRule(ctx, subscription.Name, "$Default")
			}
		}
	}
	return nil
}

// DeleteSubscription Deletes a subscription from a topic in the service bus
func (s *ServiceBusCli) DeleteSubscription(topicName string, subscriptionName string) error {
	var commonError error
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	logger.LogHighlight("Removing subscription %v from topic %v in service bus %v", log.Info, subscriptionName, topicName, s.Namespace.Name)
	topic := s.GetTopic(topicName)
	if topic == nil {
		commonError = errors.New("Could not find topic " + topicName + " in service bus" + s.Namespace.Name)
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Error, topicName, s.Namespace.Name)
		return commonError
	}
	sm := topic.NewSubscriptionManager()
	err := sm.Delete(ctx, subscriptionName)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.LogHighlight("Subscription %v was removed successfully from topic %v in service bus %v", log.Info, subscriptionName, topicName, s.Namespace.Name)
	return nil
}

// SubscribeToTopic Subscribes to a topic and listen to the messages
func (s *ServiceBusCli) SubscribeToTopic(topicName string, subscriptionName string) error {
	var commonError error

	var concurrentHandler servicebus.HandlerFunc = func(ctx context.Context, msg *servicebus.Message) error {
		logger.LogHighlight("%v Received message %v from topic %v on subscription %v with label %v", log.Info, msg.SystemProperties.EnqueuedTime.String(), msg.ID, topicName, subscriptionName, msg.Label)
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

	logger.LogHighlight("Subscribing to %v on topic %v in service bus %v", log.Info, subscriptionName, topicName, s.Namespace.Name)
	if topicName == "" {
		commonError = errors.New("Topic " + topicName + " cannot be null")
		logger.LogHighlight("Topic %v cannot be null", log.Info, topicName)
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topic := s.GetTopic(topicName)
	if topic == nil {
		commonError = errors.New("Could not find topic " + topicName + " in service bus" + s.Namespace.Name)
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Info, topicName, s.Namespace.Name)
		return commonError
	}
	s.ActiveTopic = topic

	foundSubscription := false
	sm := topic.NewSubscriptionManager()
	subscriptions, subscriptionsErr := sm.List(ctx)
	if subscriptionsErr != nil {
		logger.LogHighlight("There was an error getting the list of subscriptions on %v in service bus %v", log.Warning, topicName, s.Namespace.Name)
	}

	for _, subscription := range subscriptions {
		if subscription.Name == subscriptionName {
			foundSubscription = true
			break
		}
	}

	if !foundSubscription {
		if subscriptionName == "wiretap" {
			s.DeleteWiretap = true
			logger.LogHighlight("Wiretap subscription not found on %v in service bus %v, creating...", log.Warning, topicName, s.Namespace.Name)
			wiretapSubscription := NewSubscription(topicName, "wiretap")
			err := s.CreateSubscription(wiretapSubscription)
			if err != nil {
				logger.Error(err.Error())
				return err
			}

		} else {
			commonError := errors.New("Subscription " + subscriptionName + " was not found on " + topicName + " in service bus" + s.Namespace.Name)
			logger.LogHighlight("Subscription %v was not found on %v in service bus %v", log.Error, subscriptionName, topicName, s.Namespace.Name)
			return commonError
		}
	}

	subscription, err := topic.NewSubscription(subscriptionName)
	s.ActiveSubscription = subscription

	if err != nil {
		commonError := errors.New("Subscription " + subscriptionName + " was not found on topic " + topicName + " in service bus" + s.Namespace.Name)
		logger.LogHighlight("Subscription %v was not found on %v in service bus %v", log.Error, subscriptionName, topicName, s.Namespace.Name)
		return commonError
	}

	logger.LogHighlight("Starting to receive messages in %v on topic %v for service bus %v", log.Info, subscriptionName, topicName, s.Namespace.Name)
	receiver, err := subscription.NewReceiver(ctx)

	if err != nil {
		commonError := errors.New("Could not create channel for subscription " + subscriptionName + " on " + topicName + " in " + s.Namespace.Name + " bus, subscription was not found")
		logger.LogHighlight("Could not create channel for subscription %v on topic %v for service bus %v, subscription was not found", log.Info, subscriptionName, topicName, s.Namespace.Name)
		return commonError
	}

	listenerHandler := receiver.Listen(ctx, concurrentHandler)
	s.ActiveTopicListenerHandle = listenerHandler
	defer listenerHandler.Close(ctx)

	if <-s.CloseTopicListener {
		s.CloseTopicSubscription()
	}
	return nil
}

// CloseTopicSubscription closes the subscription to a topic
func (s *ServiceBusCli) CloseTopicSubscription() error {
	logger.LogHighlight("Closing the subscription for %v on topic %v in service bus %v", log.Info, s.ActiveSubscription.Name, s.ActiveTopic.Name, s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	s.ActiveTopicListenerHandle.Close(ctx)
	if s.DeleteWiretap && s.ActiveSubscription.Name == "wiretap" {
		s.DeleteSubscription(s.ActiveTopic.Name, "wiretap")
	}
	s.ActiveTopic = nil
	s.ActiveTopicListenerHandle = nil
	s.ActiveSubscription = nil
	s.CloseTopicListener <- false
	return nil
}
