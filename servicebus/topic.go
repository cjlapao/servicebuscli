package servicebus

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/servicebuscli-go/entities"
)

// GetTopicManager creates a servicebus topic manager
func (s *ServiceBusCli) GetTopicManager() *servicebus.TopicManager {
	logger.Trace("Creating a service bus topic manager")
	if s.Namespace == nil {
		_, err := s.GetNamespace()
		if err != nil {
			return nil
		}
	}

	s.TopicManager = s.Namespace.NewTopicManager()
	return s.TopicManager
}

// GetTopic Gets a topic from the service bus
func (s *ServiceBusCli) GetTopic(name string) *servicebus.Topic {
	logger.Trace("Getting a topic " + name + " entity in service bus " + s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if name == "" {
		return nil
	}

	if s.TopicManager == nil {
		s.GetTopicManager()
	}

	te, err := s.TopicManager.Get(ctx, name)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	topic, err := s.Namespace.NewTopic(te.Name)

	return topic
}

// GetTopicDetails Gets a topic details from the service bus
func (s *ServiceBusCli) GetTopicDetails(name string) *servicebus.TopicEntity {
	logger.Trace("Getting a topic " + name + " entity in service bus " + s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if name == "" {
		return nil
	}

	if s.TopicManager == nil {
		s.GetTopicManager()
	}

	te, err := s.TopicManager.Get(ctx, name)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return te
}

// ListTopics Lists all the topics in a service bus
func (s *ServiceBusCli) ListTopics() ([]*servicebus.TopicEntity, error) {
	logger.LogHighlight("Getting all topics in %v service bus", log.Info, s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	if s.TopicManager == nil {
		tm := s.GetTopicManager()
		if tm == nil {
			return nil, errors.New("There was an error getting the topic manager")
		}
	}

	return s.TopicManager.List(ctx)
}

// SendTopicMessage sends a message to a specific topic
func (s *ServiceBusCli) SendTopicMessage(topicName string, message entities.MessageRequest) error {
	var commonError error
	logger.LogHighlight("Sending a service bus topic message to %v topic in service bus %v", log.Info, topicName, s.Namespace.Name)
	if topicName == "" {
		commonError = errors.New("Topic cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topic := s.GetTopic(topicName)
	if topic == nil {
		commonError = errors.New("Could not find topic " + topicName + " in service bus " + s.Namespace.Name)
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Error, topicName, s.Namespace.Name)
		return commonError
	}

	sbMessage, err := message.ToServiceBus()

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = topic.Send(ctx, sbMessage)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.LogHighlight("Service bus topic message was sent successfully to %v topic in service bus %v", log.Info, topicName, s.Namespace.Name)
	logger.Info("Message Body:")
	logger.Info(string(sbMessage.Data))
	return nil
}

func (s *ServiceBusCli) SendParallelBulkTopicMessage(wg *sync.WaitGroup, topicName string, messages ...entities.MessageRequest) {
	defer wg.Done()
	_ = s.SendBulkTopicMessage(topicName, messages...)
}

func (s *ServiceBusCli) SendBulkTopicMessage(topicName string, messages ...entities.MessageRequest) error {
	var commonError error
	logger.LogHighlight("Sending a service bus topic messages to %v topic in service bus %v", log.Info, topicName, s.Namespace.Name)
	if topicName == "" {
		commonError = errors.New("Topic cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topic := s.GetTopic(topicName)
	if topic == nil {
		commonError = errors.New("Could not find topic " + topicName + " in service bus " + s.Namespace.Name)
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Error, topicName, s.Namespace.Name)
		return commonError
	}

	sbMessages := make([]*servicebus.Message, 0)
	for _, msg := range messages {
		sbMessage, err := msg.ToServiceBus()

		if err != nil {
			logger.Error(err.Error())
			return err
		}
		sbMessages = append(sbMessages, sbMessage)
	}

	err := topic.SendBatch(ctx, servicebus.NewMessageBatchIterator(262144, sbMessages...))

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.LogHighlight("Service bus bulk topic messages were sent successfully to %v topic in service bus %v", log.Info, topicName, s.Namespace.Name)
	return nil
}

// SendTopicMessage sends a message to a specific topic
func (s *ServiceBusCli) SendTopicServiceBusMessage(topicName string, sbMessage *servicebus.Message) error {
	var commonError error
	logger.LogHighlight("Sending a service bus topic message to %v topic in service bus %v", log.Info, topicName, s.Namespace.Name)
	if topicName == "" {
		commonError = errors.New("Topic cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topic := s.GetTopic(topicName)
	if topic == nil {
		commonError = errors.New("Could not find topic " + topicName + " in service bus " + s.Namespace.Name)
		logger.LogHighlight("Could not find topic %v in service bus %v", log.Error, topicName, s.Namespace.Name)
		return commonError
	}

	err := topic.Send(ctx, sbMessage)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.LogHighlight("Service bus topic message was sent successfully to %v topic in service bus %v", log.Info, topicName, s.Namespace.Name)
	logger.Info("Message:")
	logger.Info(string(sbMessage.Data))
	return nil
}

// CreateTopic Creates a topic in the service bus namespace
func (s *ServiceBusCli) CreateTopic(topicName string, opts ...servicebus.TopicManagementOption) (*servicebus.TopicEntity, error) {
	var commonError error
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if topicName == "" {
		commonError = errors.New("topic cannot be null")
		logger.Error(commonError.Error())
		return nil, commonError
	}

	logger.LogHighlight("Creating topic %v in service bus %v", log.Info, topicName, s.Namespace.Name)
	tm := s.GetTopicManager()
	topic, err := tm.Put(ctx, topicName, opts...)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	logger.LogHighlight("Topic %v was created successfully in service bus %v", log.Info, topicName, s.Namespace.Name)
	return topic, nil
}

// DeleteTopic Deletes a topic in the service bus namespace
func (s *ServiceBusCli) DeleteTopic(topicName string) error {
	var commonError error
	if topicName == "" {
		commonError = errors.New("topic cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	logger.LogHighlight("Removing topic %v in service bus %v", log.Info, topicName, s.Namespace.Name)
	tm := s.GetTopicManager()

	err := tm.Delete(ctx, topicName)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.LogHighlight("Topic %v was removed successfully from service bus %v", log.Info, topicName, s.Namespace.Name)
	return nil
}
