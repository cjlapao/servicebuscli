package entities

import (
	"time"

	azservicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/common-go/duration"
)

type SubscriptionResponse struct {
	Name                                      string            `json:"name"`
	ID                                        string            `json:"url"`
	CountDetails                              CountDetails      `json:"countDetails"`
	LockDuration                              duration.Duration `json:"lockDuration"`
	RequiresSession                           *bool             `json:"requiresSession"`
	DefaultMessageTimeToLive                  duration.Duration `json:"defaultMessageTimeToLive"`
	AutoDeleteOnIdle                          duration.Duration `json:"autoDeleteOnIdle"`
	DeadLetteringOnMessageExpiration          *bool             `json:"deadLetteringOnMessageExpiration"`
	DeadLetteringOnFilterEvaluationExceptions *bool             `json:"deadLetteringOnFilterEvaluationExceptions"`
	MessageCount                              *int64            `json:"messageCount"`
	MaxDeliveryCount                          int64             `json:"maxDeliveryCount"`
	EnableBatchedOperations                   *bool             `json:"enableBatchedOperations"`
	Status                                    string            `json:"status"`
	CreatedAt                                 time.Time         `json:"createdAt"`
	UpdatedAt                                 time.Time         `json:"updatedAt"`
	AccessedAt                                time.Time         `json:"accessedAt"`
	ForwardTo                                 *string           `json:"forwardTo"`
	ForwardDeadLetteredMessagesTo             *string           `json:"forwardDeadLetteredMessagesTo"`
}

func (e *SubscriptionResponse) FromServiceBus(subscription *azservicebus.SubscriptionEntity) {
	lockDuration, _ := duration.FromString(*subscription.LockDuration)
	defaultMessageTimeToLive, _ := duration.FromString(*subscription.DefaultMessageTimeToLive)
	autoDeleteOnIdle, _ := duration.FromString(*subscription.AutoDeleteOnIdle)

	e.LockDuration = *lockDuration
	e.RequiresSession = subscription.RequiresSession
	e.AutoDeleteOnIdle = *autoDeleteOnIdle
	e.DefaultMessageTimeToLive = *defaultMessageTimeToLive
	e.DeadLetteringOnMessageExpiration = subscription.DeadLetteringOnMessageExpiration
	e.DeadLetteringOnFilterEvaluationExceptions = subscription.DeadLetteringOnFilterEvaluationExceptions
	e.MessageCount = subscription.MessageCount
	e.MaxDeliveryCount = int64(*subscription.MaxDeliveryCount)
	e.EnableBatchedOperations = subscription.EnableBatchedOperations
	e.Status = string(*subscription.Status)
	e.CreatedAt = subscription.CreatedAt.Time
	e.UpdatedAt = subscription.UpdatedAt.Time
	e.AccessedAt = subscription.AccessedAt.Time
	e.ForwardTo = subscription.ForwardTo
	e.ForwardDeadLetteredMessagesTo = subscription.ForwardDeadLetteredMessagesTo
	e.Name = subscription.Name
	e.ID = subscription.ID

	e.CountDetails = CountDetails{
		ActiveMessageCount:             subscription.CountDetails.ActiveMessageCount,
		DeadLetterMessageCount:         subscription.CountDetails.DeadLetterMessageCount,
		ScheduledMessageCount:          subscription.CountDetails.ScheduledMessageCount,
		TransferDeadLetterMessageCount: subscription.CountDetails.TransferDeadLetterMessageCount,
		TransferMessageCount:           subscription.CountDetails.TransferMessageCount,
	}
}
