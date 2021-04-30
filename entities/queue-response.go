package entities

import (
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/servicebuscli-go/duration"
)

// QueueResponse
type QueueResponse struct {
	Name                                string             `json:"name"`
	ID                                  string             `json:"id"`
	CountDetails                        CountDetails       `json:"countDetails"`
	LockDuration                        *duration.Duration `json:"lockDuration"`
	MaxSizeInMegabytes                  *int32             `json:"maxSizeInMegabytes"`
	RequiresDuplicateDetection          *bool              `json:"requiresDuplicateDetection"`
	RequiresSession                     *bool              `json:"requiresSession"`
	DefaultMessageTimeToLive            *duration.Duration `json:"defaultMessageTimeToLive"`
	DeadLetteringOnMessageExpiration    *bool              `json:"deadLetteringOnMessageExpiration"`
	DuplicateDetectionHistoryTimeWindow *duration.Duration `json:"duplicateDetectionHistoryTimeWindow"`
	MaxDeliveryCount                    *int32             `json:"maxDeliveryCount"`
	EnableBatchedOperations             *bool              `json:"enableBatchedOperations"`
	SizeInBytes                         *int64             `json:"sizeInBytes"`
	MessageCount                        *int64             `json:"messageCount"`
	IsAnonymousAccessible               *bool              `json:"isAnonymousAccessible"`
	Status                              string             `json:"status"`
	CreatedAt                           time.Time          `json:"createdAt"`
	UpdatedAt                           time.Time          `json:"updatedAt"`
	SupportOrdering                     *bool              `json:"supportOrdering"`
	AutoDeleteOnIdle                    *duration.Duration `json:"autoDeleteOnIdle"`
	EnablePartitioning                  *bool              `json:"enablePartitioning"`
	EnableExpress                       *bool              `json:"enableExpress"`
	ForwardTo                           *string            `json:"forwardTo"`
	ForwardDeadLetteredMessagesTo       *string            `json:"forwardDeadLetteredMessagesTo"`
}

func (q *QueueResponse) FromServiceBus(queue *servicebus.QueueEntity) {
	if queue == nil {
		return
	}

	q.Name = queue.Name
	q.ID = queue.ID
	lockDuration, _ := duration.FromString(*queue.LockDuration)
	defaultMessageTimeToLive, _ := duration.FromString(*queue.DefaultMessageTimeToLive)
	duplicateDetectionHistoryTimeWindow, _ := duration.FromString(*queue.DuplicateDetectionHistoryTimeWindow)
	autoDeleteOnIdle, _ := duration.FromString(*queue.AutoDeleteOnIdle)
	q.LockDuration = lockDuration
	q.DefaultMessageTimeToLive = defaultMessageTimeToLive
	q.MaxSizeInMegabytes = queue.MaxSizeInMegabytes
	q.RequiresDuplicateDetection = queue.RequiresDuplicateDetection
	q.RequiresSession = queue.RequiresSession
	q.DeadLetteringOnMessageExpiration = queue.DeadLetteringOnMessageExpiration
	q.DuplicateDetectionHistoryTimeWindow = duplicateDetectionHistoryTimeWindow
	q.MaxDeliveryCount = queue.MaxDeliveryCount
	q.EnableBatchedOperations = queue.EnableBatchedOperations
	q.SizeInBytes = queue.SizeInBytes
	q.IsAnonymousAccessible = queue.IsAnonymousAccessible
	q.Status = string(*queue.Status)
	q.CreatedAt = queue.CreatedAt.Time
	q.UpdatedAt = queue.UpdatedAt.Time
	q.SupportOrdering = queue.SupportOrdering
	q.AutoDeleteOnIdle = autoDeleteOnIdle
	q.EnablePartitioning = queue.EnablePartitioning
	q.EnableExpress = queue.EnableExpress
	q.ForwardTo = queue.ForwardTo
	q.ForwardDeadLetteredMessagesTo = queue.ForwardDeadLetteredMessagesTo
	q.CountDetails = CountDetails{}
	q.CountDetails.FromServiceBus(queue.CountDetails)
}
