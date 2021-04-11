package entities

import (
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type TopicEntity struct {
	Name                                string       `json:"name"`
	ID                                  string       `json:"id"`
	CountDetails                        CountDetails `json:"countDetails,omitempty"`
	DefaultMessageTimeToLive            *string      `json:"defaultMessageTimeToLive"`
	MaxSizeInMegabytes                  *int32       `json:"maxSizeInMegabytes"`
	RequiresDuplicateDetection          *bool        `json:"requiresDuplicateDetection"`
	DuplicateDetectionHistoryTimeWindow *string      `json:"duplicateDetectionHistoryTimeWindow"`
	EnableBatchedOperations             *bool        `json:"enableBatchedOperations"`
	SizeInBytes                         *int64       `json:"sizeInBytes"`
	FilteringMessagesBeforePublishing   *bool        `json:"filteringMessagesBeforePublishing"`
	IsAnonymousAccessible               *bool        `json:"isAnonymousAccessible"`
	Status                              string       `json:"status"`
	CreatedAt                           time.Time    `json:"createdAt"`
	UpdatedAt                           time.Time    `json:"updatedAt"`
	SupportOrdering                     *bool        `json:"supportOrdering"`
	AutoDeleteOnIdle                    *string      `json:"autoDeleteOnIdle"`
	EnablePartitioning                  *bool        `json:"enablePartitioning"`
	EnableSubscriptionPartitioning      *bool        `json:"enableSubscriptionPartitioning"`
	EnableExpress                       *bool        `json:"enableExpress"`
}

func (t *TopicEntity) FromServiceBus(topic *servicebus.TopicEntity) {
	if topic == nil {
		return
	}

	t.Name = topic.Name
	t.ID = topic.ID
	t.DefaultMessageTimeToLive = topic.DefaultMessageTimeToLive
	t.MaxSizeInMegabytes = topic.MaxSizeInMegabytes
	t.RequiresDuplicateDetection = topic.RequiresDuplicateDetection
	t.DuplicateDetectionHistoryTimeWindow = topic.DuplicateDetectionHistoryTimeWindow
	t.EnableBatchedOperations = topic.EnableBatchedOperations
	t.SizeInBytes = topic.SizeInBytes
	t.FilteringMessagesBeforePublishing = topic.FilteringMessagesBeforePublishing
	t.IsAnonymousAccessible = topic.IsAnonymousAccessible
	t.Status = string(*topic.Status)
	t.CreatedAt = topic.CreatedAt.Time
	t.UpdatedAt = topic.UpdatedAt.Time
	t.SupportOrdering = topic.SupportOrdering
	t.AutoDeleteOnIdle = topic.AutoDeleteOnIdle
	t.EnablePartitioning = topic.EnablePartitioning
	t.EnableSubscriptionPartitioning = topic.EnableSubscriptionPartitioning
	t.EnableExpress = topic.EnableExpress
	t.CountDetails = CountDetails{}
	t.CountDetails.FromServiceBus(topic.CountDetails)
}

type TopicRequest struct {
	Name    string               `json:"name"`
	Options *TopicRequestOptions `json:"options,omitempty"`
}

type TopicRequestOptions struct {
	AutoDeleteOnIdle         *string `json:"autoDeleteOnIdle,omitempty"`
	EnableBatchedOperation   *bool   `json:"enableBatchedOperation,omitempty"`
	EnableDuplicateDetection *bool   `json:"enableDuplicateDetection,omitempty"`
	EnableExpress            *bool   `json:"enableExpress,omitempty"`
	MaxSizeInMegabytes       *int    `json:"maxSizeInMegabytes,omitempty"`
	DefaultMessageTimeToLive *string `json:"defaultMessageTimeToLive,omitempty"`
	SupportOrdering          *bool   `json:"supportOrdering,omitempty"`
	EnablePartitioning       *bool   `json:"enablePartitioning,omitempty"`
}
