package servicebus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/servicebuscli-go/entities"
)

// GetSubscription Gets a subscription from a topic in the namespace
func (s *ServiceBusCli) GetSubscription(topicName string, subscriptionName string) (*servicebus.SubscriptionEntity, error) {
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
	return sm.Get(ctx, subscriptionName)
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
func (s *ServiceBusCli) CreateSubscription(subscription entities.SubscriptionRequestEntity, upsert bool) error {
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
	existingSubscription, _ := sm.Get(ctx, subscription.Name)
	if existingSubscription != nil && !upsert {
		commonError = errors.New("Subscription " + subscription.Name + " already exists on topic " + subscription.TopicName + " in service bus" + s.Namespace.Name)
		logger.LogHighlight("Subscription %v already exists on topic %v in service bus %v", log.Error, subscription.Name, subscription.TopicName, s.Namespace.Name)
		return commonError
	}

	// Generating subscription options
	entityOpts, apiErr := subscription.GetOptions()
	if apiErr != nil {
		logger.Error("There was an error creating subscription")
		logger.Error(apiErr.Message)
		return errors.New(apiErr.Message)
	}

	opts = append(opts, *entityOpts...)

	// Generating the forward rule, checking if the target exists or not
	if subscription.Forward != nil && subscription.Forward.To != "" {
		switch subscription.Forward.In {
		case entities.ForwardToTopic:
			tm := s.GetTopicManager()
			target, err := tm.Get(ctx, subscription.Forward.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding topic %v in service bus %v", log.Error, subscription.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.SubscriptionWithAutoForward(target))
		case entities.ForwardToQueue:
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
	if subscription.ForwardDeadLetter != nil && subscription.ForwardDeadLetter.To != "" {
		switch subscription.ForwardDeadLetter.In {
		case entities.ForwardToTopic:
			tm := s.GetTopicManager()
			target, err := tm.Get(ctx, subscription.ForwardDeadLetter.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding topic %v in service bus %v", log.Error, subscription.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.SubscriptionWithAutoForward(target))
		case entities.ForwardToQueue:
			qm := s.GetQueueManager()
			target, err := qm.Get(ctx, subscription.ForwardDeadLetter.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding queue %v in service bus %v", log.Error, subscription.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.SubscriptionWithForwardDeadLetteredMessagesTo(target))
		}
	}

	_, err := sm.Put(ctx, subscription.Name, opts...)
	if err != nil {
		logger.Error("There was an error creating subscription")
		logger.Error(err.Error())
		return err
	}

	// Defining the filters if they exist
	if subscription.Rules != nil {
		for _, rule := range subscription.Rules {
			s.CreateSubscriptionRule(subscription, *rule)
			if err != nil {
				return err
			}
		}
	}

	logger.LogHighlight("Subscription %v was created successfully on topic %v in service bus %v", log.Info, subscription.Name, subscription.TopicName, s.Namespace.Name)
	return nil
}

func (s *ServiceBusCli) GetSubscriptionRules(topicName string, subscriptionName string) ([]*servicebus.RuleEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	logger.LogHighlight("Getting subscription rules in subscription %v on topic %v in service bus %v", log.Info, subscriptionName, topicName, s.Namespace.Name)
	topic := s.GetTopic(topicName)
	sm := topic.NewSubscriptionManager()

	return sm.ListRules(ctx, subscriptionName)
}

func (s *ServiceBusCli) GetSubscriptionRule(topicName string, subscriptionName string, ruleName string) (*servicebus.RuleEntity, error) {
	rules, err := s.GetSubscriptionRules(topicName, subscriptionName)
	if err != nil {
		return nil, err
	}

	logger.LogHighlight("Trying to find the rule %v, in subscription %v on topic %v in service bus %v", log.Info, ruleName, subscriptionName, topicName, s.Namespace.Name)
	for _, rule := range rules {
		if strings.ToLower(rule.Name) == strings.ToLower(ruleName) {
			return rule, nil
		}
	}

	notFound := errors.New("No rule was found with name " + ruleName)
	return nil, notFound
}

func (s *ServiceBusCli) DeleteSubscriptionRule(topicName string, subscriptionName string, ruleName string) (*servicebus.RuleEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	logger.LogHighlight("Deleting subscription rule %v in subscription %v on topic %v in service bus %v", log.Info, ruleName, subscriptionName, topicName, s.Namespace.Name)
	topic := s.GetTopic(topicName)
	sm := topic.NewSubscriptionManager()

	rule, err := s.GetSubscriptionRule(topicName, subscriptionName, ruleName)

	if err != nil {
		return nil, err
	}

	err = sm.DeleteRule(ctx, subscriptionName, ruleName)

	if err != nil {
		return nil, err
	}

	return rule, nil
}

// CreateSubscriptionRule Creates a rule to a specific subscription
func (s *ServiceBusCli) CreateSubscriptionRule(subscription entities.SubscriptionRequestEntity, rule entities.RuleRequestEntity) error {
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
		logger.LogHighlight("Topic %v cannot be null", log.Error, topicName)
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topic := s.GetTopic(topicName)
	if topic == nil {
		commonError = errors.New("Could not find topic " + topicName + " in service bus" + s.Namespace.Name)
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Error, topicName, s.Namespace.Name)
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
			logger.LogHighlight("Wiretap subscription not found on %v in service bus %v, creating...", log.Info, topicName, s.Namespace.Name)
			wiretapSubscription := entities.NewSubscriptionRequest(topicName, "wiretap")
			err := s.CreateSubscription(*wiretapSubscription, false)
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
		logger.LogHighlight("Could not create channel for subscription %v on topic %v for service bus %v, subscription was not found", log.Error, subscriptionName, topicName, s.Namespace.Name)
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

// GetSubscriptionActiveMessages Gets messages from a subscription
func (s *ServiceBusCli) GetSubscriptionActiveMessages(topicName string, subscriptionName string, qty int, peek bool) ([]servicebus.Message, error) {
	var commonError error
	messages := make([]servicebus.Message, 0)

	// We will have a maximum of fetch of 100 messages per query
	if qty > 100 {
		qty = 100
	}

	logger.LogHighlight("Getting message for subscription %v on topic %v in service bus %v", log.Info, subscriptionName, topicName, s.Namespace.Name)
	if topicName == "" {
		commonError = errors.New("Topic " + topicName + " cannot be null")
		logger.LogHighlight("Topic %v cannot be null", log.Error, topicName)
		return messages, commonError
	}

	topic := s.GetTopic(topicName)
	if topic == nil {
		commonError = errors.New("Could not find topic " + topicName + " in service bus namespace" + s.Namespace.Name)
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Error, topicName, s.Namespace.Name)
		return messages, commonError
	}

	s.ActiveTopic = topic
	sm := topic.NewSubscriptionManager()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	subscription, commonError := sm.Get(ctx, subscriptionName)
	if commonError != nil {
		return messages, commonError
	}

	if *subscription.CountDetails.ActiveMessageCount <= 0 {
		return messages, nil
	}

	messageCount := int(*subscription.CountDetails.ActiveMessageCount)

	// If we set the message count to 0 then we will read a batch of the messages that exists
	if qty == 0 {
		qty = messageCount
		// we will need to check again that the message count is not bigger than 100 and will set the limit again if so
		if qty > 100 {
			qty = 100
		}
	}

	// if we do not have as many messages in the system as requested we will adjust the quantity to the amount of messages in the system
	if messageCount < qty {
		qty = messageCount
	}

	var waitForMessages sync.WaitGroup
	waitForMessages.Add(qty)

	// Creating a message handler function to deal with our messages
	var messageHandler servicebus.HandlerFunc = func(c context.Context, m *servicebus.Message) error {
		mCtx, mCancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
		defer mCancel()
		m.Complete(mCtx)

		messages = append(messages, *m)
		// notify the wait group that the message is dealt with
		waitForMessages.Done()
		return nil
	}

	// Creating the receiver for the messages in the subscription
	messageReceiver, commonError := topic.NewSubscription(subscriptionName)

	if commonError != nil {
		return nil, commonError
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Background task to receive all of the messages we need
	go func() {
		for i := 0; i < qty; i++ {
			if peek {
				m, commonError := messageReceiver.PeekOne(ctx)
				if commonError == nil {
					messages = append(messages, *m)
				}
				waitForMessages.Done()
			} else {
				if err := messageReceiver.ReceiveOne(ctx, messageHandler); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}()

	waitForMessages.Wait()

	// We are finished and we should now close the receiver before leaving
	_ = topic.Close(ctx)
	_ = messageReceiver.Close(ctx)

	return messages, nil
}

func (s *ServiceBusCli) GetSubscriptionDeadLetterMessages(topicName string, subscriptionName string, qty int, peek bool) ([]servicebus.Message, error) {
	var commonError error
	messages := make([]servicebus.Message, 0)

	// We will have a maximum of fetch of 100 messages per query
	if qty > 100 {
		qty = 100
	}

	logger.LogHighlight("Getting dead letter messages for subscription %v on topic %v in service bus %v", log.Info, subscriptionName, topicName, s.Namespace.Name)
	if topicName == "" {
		commonError = errors.New("Topic " + topicName + " cannot be null")
		logger.LogHighlight("Topic %v cannot be null", log.Error, topicName)
		return messages, commonError
	}

	topic := s.GetTopic(topicName)
	if topic == nil {
		commonError = errors.New("Could not find topic " + topicName + " in service bus namespace" + s.Namespace.Name)
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Error, topicName, s.Namespace.Name)
		return messages, commonError
	}

	s.ActiveTopic = topic
	sm := topic.NewSubscriptionManager()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	subscription, commonError := sm.Get(ctx, subscriptionName)
	if commonError != nil {
		logger.LogHighlight("Could not find subscription %v on topic %v in service bus %v", log.Error, subscriptionName, topicName, s.Namespace.Name)
		return messages, commonError
	}

	if *subscription.CountDetails.DeadLetterMessageCount <= 0 {
		return messages, nil
	}

	messageCount := int(*subscription.CountDetails.DeadLetterMessageCount)

	// If we set the message count to 0 then we will read a batch of the messages that exists
	if qty == 0 {
		qty = messageCount
		// we will need to check again that the message count is not bigger than 100 and will set the limit again if so
		if qty > 100 {
			qty = 100
		}
	}

	// if we do not have as many messages in the system as requested we will adjust the quantity to the amount of messages in the system
	if messageCount < qty {
		qty = messageCount
	}

	var waitForMessages sync.WaitGroup
	waitForMessages.Add(qty)

	// Creating a message handler function to deal with our messages
	var messageHandler servicebus.HandlerFunc = func(c context.Context, m *servicebus.Message) error {
		mCtx, mCancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
		defer mCancel()

		if peek {
			m.Abandon(ctx)
		} else {
			m.Complete(mCtx)
		}

		messages = append(messages, *m)
		// notify the wait group that the message is dealt with
		waitForMessages.Done()
		return nil
	}

	// Creating the receiver for the messages in the subscription
	messageReceiver, commonError := topic.NewSubscription(subscriptionName)

	if commonError != nil {
		return nil, commonError
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var deadLetterReceiver servicebus.ReceiveOner
	if peek {
		deadLetterReceiver, commonError = messageReceiver.NewDeadLetterReceiver(ctx, servicebus.ReceiverWithReceiveMode(servicebus.PeekLockMode))

	} else {
		deadLetterReceiver, commonError = messageReceiver.NewDeadLetterReceiver(ctx)
	}

	if commonError != nil {
		return nil, commonError
	}

	// Background task to receive all of the messages we need
	go func() {
		for i := 0; i < qty; i++ {
			if err := deadLetterReceiver.ReceiveOne(ctx, messageHandler); err != nil {
				fmt.Println(err.Error())
			}
		}
	}()

	waitForMessages.Wait()

	// We are finished and we should now close the receiver before leaving
	_ = topic.Close(ctx)
	_ = messageReceiver.Close(ctx)
	_ = deadLetterReceiver.Close(ctx)

	return messages, nil
}
