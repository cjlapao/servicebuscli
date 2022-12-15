package servicebus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/servicebuscli/entities"
)

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
func (s *ServiceBusCli) GetQueue(queueName string) (*servicebus.Queue, error) {
	logger.Trace("Getting queue " + queueName + " from service bus " + s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if queueName == "" {
		err := errors.New("queue name is empty or nil")
		return nil, err
	}

	if s.QueueManager == nil {
		s.GetQueueManager()
	}

	qe, err := s.QueueManager.Get(ctx, queueName)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return s.Namespace.NewQueue(qe.Name)
}

// GetQueueDetails Gets a Namespace Queue Entity with details
func (s *ServiceBusCli) GetQueueDetails(queueName string) (*servicebus.QueueEntity, error) {
	logger.Trace("Getting queue " + queueName + " from service bus " + s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if queueName == "" {
		err := errors.New("queue name is nil or empty")
		return nil, err
	}
	if s.QueueManager == nil {
		s.GetQueueManager()
	}

	return s.QueueManager.Get(ctx, queueName)
}

// ListQueues Lists all the Queues in a Service Bus
func (s *ServiceBusCli) ListQueues() ([]*servicebus.QueueEntity, error) {
	logger.LogHighlight("Getting all queues from %v service bus ", log.Info, s.Namespace.Name)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	qm := s.GetQueueManager()
	if qm == nil {
		commonError := errors.New("there was an error getting the queue manager, check your internet connection")
		logger.LogHighlight("There was an error getting the %v, check your internet connection", log.Error, "queue manager")
		return nil, commonError
	}

	return qm.List(ctx)
}

// CreateQueue Creates a queue in the service bus namespace
func (s *ServiceBusCli) CreateQueue(queue entities.QueueRequest) error {
	var commonError error
	opts := make([]servicebus.QueueManagementOption, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if queue.Name == "" {
		commonError = errors.New("queue name cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}
	logger.LogHighlight("Creating queue %v in service bus %v", log.Info, queue.Name, s.Namespace.Name)

	qm := s.GetQueueManager()

	// Checking if the queue already exists in the namespace
	existingQueue, _ := qm.Get(ctx, queue.Name)
	if existingQueue != nil {
		commonError = errors.New("queue " + queue.Name + " already exists in service bus " + s.Namespace.Name)
		logger.LogHighlight("Queue %v already exists in service bus %v", log.Error, queue.Name, s.Namespace.Name)
		return commonError
	}

	// Generating subscription options
	entityOpts, apiErr := queue.GetOptions()
	if apiErr != nil {
		logger.Error("There was an error creating queue")
		logger.Error(apiErr.Message)
		return errors.New(apiErr.Message)
	}
	opts = append(opts, *entityOpts...)

	if queue.MaxDeliveryCount > 0 && queue.MaxDeliveryCount != 10 {
		opts = append(opts, servicebus.QueueEntityWithMaxDeliveryCount(int32(queue.MaxDeliveryCount)))
	}

	// Generating the forward rule, checking if the target exists or not
	if queue.Forward != nil && queue.Forward.To != "" {
		switch queue.Forward.In {
		case entities.ForwardToTopic:
			tm := s.GetTopicManager()
			target, err := tm.Get(ctx, queue.Forward.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding topic %v in service bus %v", log.Error, queue.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.QueueEntityWithAutoForward(target))
		case entities.ForwardToQueue:
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
	if queue.ForwardDeadLetter != nil && queue.ForwardDeadLetter.To != "" {
		switch queue.ForwardDeadLetter.In {
		case entities.ForwardToTopic:
			tm := s.GetTopicManager()
			target, err := tm.Get(ctx, queue.ForwardDeadLetter.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding topic %v in service bus %v", log.Error, queue.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.QueueEntityWithForwardDeadLetteredMessagesTo(target))
		case entities.ForwardToQueue:
			qm := s.GetQueueManager()
			target, err := qm.Get(ctx, queue.ForwardDeadLetter.To)
			if err != nil || target == nil {
				logger.LogHighlight("Could not find forwarding queue %v in service bus %v", log.Error, queue.Forward.To, s.Namespace.Name)
				return err
			}
			opts = append(opts, servicebus.QueueEntityWithForwardDeadLetteredMessagesTo(target))
		}
	}

	_, err := qm.Put(ctx, queue.Name, opts...)
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
		commonError = errors.New("queue cannot be null")
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
func (s *ServiceBusCli) SendQueueMessage(queueName string, message entities.MessageRequest) error {
	var commonError error
	logger.LogHighlight("Sending a service bus queue message to %v queue in service bus %v", log.Info, queueName, s.Namespace.Name)
	if queueName == "" {
		commonError = errors.New("queue cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue, err := s.GetQueue(queueName)
	if queue == nil || err != nil {
		commonError = errors.New("Could not find queue " + queueName + " in service bus " + s.Namespace.Name)
		logger.LogHighlight("Could not find queue %v in service bus %v", log.Error, queueName, s.Namespace.Name)
		return commonError
	}

	messageData, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	sbMessage, err := message.ToServiceBus()

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = queue.Send(ctx, sbMessage)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.LogHighlight("Service bus queue message was sent successfully to %v queue in service bus %v", log.Info, queueName, s.Namespace.Name)
	logger.Info("Message:")
	logger.Info(string(messageData))
	return nil
}

func (s *ServiceBusCli) SendParallelBulkQueueMessage(wg *sync.WaitGroup, queueName string, messages ...entities.MessageRequest) {
	defer wg.Done()
	_ = s.SendBulkQueueMessage(queueName, messages...)
}

// SendQueueMessage Sends a Service Bus Message to a Queue
func (s *ServiceBusCli) SendBulkQueueMessage(queueName string, messages ...entities.MessageRequest) error {
	var commonError error
	logger.LogHighlight("Sending a service bus queue messages to %v queue in service bus %v", log.Info, queueName, s.Namespace.Name)
	if queueName == "" {
		commonError = errors.New("queue cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue, err := s.GetQueue(queueName)
	if queue == nil || err != nil {
		commonError = errors.New("Could not find queue " + queueName + " in service bus " + s.Namespace.Name)
		logger.LogHighlight("Could not find queue %v in service bus %v", log.Error, queueName, s.Namespace.Name)
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

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = queue.SendBatch(ctx, servicebus.NewMessageBatchIterator(262144, sbMessages...))

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.LogHighlight("Service bus bulk queue messages were sent successfully to %v queue in service bus %v", log.Info, queueName, s.Namespace.Name)
	return nil
}

// SendQueueServiceBusMessage Sends a Service Bus Message to a Queue
func (s *ServiceBusCli) SendQueueServiceBusMessage(queueName string, sbMessage *servicebus.Message) error {
	var commonError error
	logger.LogHighlight("Sending a service bus queue message to %v queue in service bus %v", log.Info, queueName, s.Namespace.Name)
	if queueName == "" {
		commonError = errors.New("queue cannot be null")
		logger.Error(commonError.Error())
		return commonError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue, err := s.GetQueue(queueName)
	if queue == nil || err != nil {
		commonError = errors.New("Could not find queue " + queueName + " in service bus " + s.Namespace.Name)
		logger.LogHighlight("Could not find queue %v in service bus %v", log.Error, queueName, s.Namespace.Name)
		return commonError
	}

	err = queue.Send(ctx, sbMessage)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.LogHighlight("Service bus queue message was sent successfully to %v queue in service bus %v", log.Info, queueName, s.Namespace.Name)
	logger.Info("Message:")
	logger.Info(string(sbMessage.Data))
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

	queue, err := s.GetQueue(queueName)
	if queue == nil || err != nil {
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

// GetSubscriptionActiveMessages Gets messages from a subscription
func (s *ServiceBusCli) GetQueueActiveMessages(queueName string, qty int, peek bool) ([]servicebus.Message, error) {
	var commonError error
	messages := make([]servicebus.Message, 0)

	// We will have a maximum of fetch of 100 messages per query
	if qty > 100 {
		qty = 100
	}

	logger.LogHighlight("Getting message for queue %v in service bus %v", log.Info, queueName, s.Namespace.Name)
	queue, _ := s.GetQueue(queueName)
	if queue == nil {
		commonError = errors.New("Could not find queue " + queueName + " in service bus namespace" + s.Namespace.Name)
		logger.LogHighlight("Could not find queue %v in service bus %v", log.Error, queueName, s.Namespace.Name)
		return messages, commonError
	}

	s.ActiveQueue = queue
	qm := s.GetQueueManager()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	queueEntity, commonError := qm.Get(ctx, queueName)
	if commonError != nil {
		return messages, commonError
	}

	if *queueEntity.CountDetails.ActiveMessageCount <= 0 {
		return messages, nil
	}

	messageCount := int(*queueEntity.CountDetails.ActiveMessageCount)

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
	messageReceiver, commonError := queue.NewReceiver(ctx)

	if commonError != nil {
		return nil, commonError
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Background task to receive all of the messages we need
	go func() {
		if peek {
			t, err := queue.Peek(ctx)
			if err != nil {
				fmt.Println(err.Error())
			}
			for i := 0; i < qty; i++ {
				m, err := t.Next(ctx)
				if err != nil {
					fmt.Println(err.Error())
				}
				messages = append(messages, *m)
				waitForMessages.Done()
			}
		} else {
			for i := 0; i < qty; i++ {
				if err := queue.ReceiveOne(ctx, messageHandler); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}()

	waitForMessages.Wait()

	// We are finished and we should now close the receiver before leaving
	_ = queue.Close(ctx)
	_ = messageReceiver.Close(ctx)

	return messages, nil
}

func (s *ServiceBusCli) GetQueueDeadLetterMessages(queueName string, qty int, peek bool) ([]servicebus.Message, error) {
	var commonError error
	messages := make([]servicebus.Message, 0)

	// We will have a maximum of fetch of 100 messages per query
	if qty > 100 {
		qty = 100
	}

	logger.LogHighlight("Getting dead letter messages for queue %v in service bus %v", log.Info, queueName, s.Namespace.Name)
	if queueName == "" {
		commonError = errors.New("Topic " + queueName + " cannot be null")
		logger.LogHighlight("Topic %v cannot be null", log.Error, queueName)
		return messages, commonError
	}

	queue, _ := s.GetQueue(queueName)
	if queue == nil {
		commonError = errors.New("Could not find queue " + queueName + " in service bus namespace" + s.Namespace.Name)
		logger.LogHighlight("Could not find queue %v in service bus %v", log.Error, queueName, s.Namespace.Name)
		return messages, commonError
	}

	s.ActiveQueue = queue
	qm := s.GetQueueManager()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	queueEntity, _ := qm.Get(ctx, queueName)

	if *queueEntity.CountDetails.DeadLetterMessageCount <= 0 {
		return messages, nil
	}

	messageCount := int(*queueEntity.CountDetails.DeadLetterMessageCount)

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
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var deadLetterReceiver servicebus.ReceiveOner
	if peek {
		deadLetterReceiver, commonError = queue.NewDeadLetterReceiver(ctx, servicebus.ReceiverWithReceiveMode(servicebus.PeekLockMode))

	} else {
		deadLetterReceiver, commonError = queue.NewDeadLetterReceiver(ctx)
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
	_ = queue.Close(ctx)
	_ = deadLetterReceiver.Close(ctx)

	return messages, nil
}
