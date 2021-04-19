package entities

import (
	"time"

	azservicebus "github.com/Azure/azure-service-bus-go"
)

type SubscriptionEntity struct {
	Name                                      string             `json:"name"`
	ID                                        string             `json:"url"`
	CountDetails                              CountDetailsEntity `json:"countDetails"`
	LockDuration                              *string            `json:"lockDuration"`
	RequiresSession                           *bool              `json:"requiresSession"`
	DefaultMessageTimeToLive                  *string            `json:"defaultMessageTimeToLive"`
	DeadLetteringOnMessageExpiration          *bool              `json:"deadLetteringOnMessageExpiration"`
	DeadLetteringOnFilterEvaluationExceptions *bool              `json:"deadLetteringOnFilterEvaluationExceptions"`
	MessageCount                              *int64             `json:"messageCount"`
	MaxDeliveryCount                          *int32             `json:"maxDeliveryCount"`
	EnableBatchedOperations                   *bool              `json:"enableBatchedOperations"`
	Status                                    string             `json:"status"`
	CreatedAt                                 time.Time          `json:"createdAt"`
	UpdatedAt                                 time.Time          `json:"updatedAt"`
	AccessedAt                                time.Time          `json:"accessedAt"`
	ForwardTo                                 *string            `json:"forwardTo"`
	ForwardDeadLetteredMessagesTo             *string            `json:"forwardDeadLetteredMessagesTo"`
}

func (e *SubscriptionEntity) FromServiceBus(subscription *azservicebus.SubscriptionEntity) {
	e.LockDuration = subscription.LockDuration
	e.RequiresSession = subscription.RequiresSession
	e.DefaultMessageTimeToLive = subscription.DefaultMessageTimeToLive
	e.DeadLetteringOnMessageExpiration = subscription.DeadLetteringOnMessageExpiration
	e.DeadLetteringOnFilterEvaluationExceptions = subscription.DeadLetteringOnFilterEvaluationExceptions
	e.MessageCount = subscription.MessageCount
	e.MaxDeliveryCount = subscription.MaxDeliveryCount
	e.EnableBatchedOperations = subscription.EnableBatchedOperations
	e.Status = string(*subscription.Status)
	e.CreatedAt = subscription.CreatedAt.Time
	e.UpdatedAt = subscription.UpdatedAt.Time
	e.AccessedAt = subscription.AccessedAt.Time
	e.ForwardTo = subscription.ForwardTo
	e.ForwardDeadLetteredMessagesTo = subscription.ForwardDeadLetteredMessagesTo
	e.Name = subscription.Name
	e.ID = subscription.ID

	e.CountDetails = CountDetailsEntity{
		ActiveMessageCount:             subscription.CountDetails.ActiveMessageCount,
		DeadLetterMessageCount:         subscription.CountDetails.DeadLetterMessageCount,
		ScheduledMessageCount:          subscription.CountDetails.ScheduledMessageCount,
		TransferDeadLetterMessageCount: subscription.CountDetails.TransferDeadLetterMessageCount,
		TransferMessageCount:           subscription.CountDetails.TransferMessageCount,
	}
}
