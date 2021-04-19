package entities

import (
	"net/http"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type TopicEntity struct {
	Name                                string             `json:"name"`
	ID                                  string             `json:"id"`
	CountDetails                        CountDetailsEntity `json:"countDetails,omitempty"`
	DefaultMessageTimeToLive            *string            `json:"defaultMessageTimeToLive"`
	MaxSizeInMegabytes                  *int32             `json:"maxSizeInMegabytes"`
	RequiresDuplicateDetection          *bool              `json:"requiresDuplicateDetection"`
	DuplicateDetectionHistoryTimeWindow *string            `json:"duplicateDetectionHistoryTimeWindow"`
	EnableBatchedOperations             *bool              `json:"enableBatchedOperations"`
	SizeInBytes                         *int64             `json:"sizeInBytes"`
	FilteringMessagesBeforePublishing   *bool              `json:"filteringMessagesBeforePublishing"`
	IsAnonymousAccessible               *bool              `json:"isAnonymousAccessible"`
	Status                              string             `json:"status"`
	CreatedAt                           time.Time          `json:"createdAt"`
	UpdatedAt                           time.Time          `json:"updatedAt"`
	SupportOrdering                     *bool              `json:"supportOrdering"`
	AutoDeleteOnIdle                    *string            `json:"autoDeleteOnIdle"`
	EnablePartitioning                  *bool              `json:"enablePartitioning"`
	EnableSubscriptionPartitioning      *bool              `json:"enableSubscriptionPartitioning"`
	EnableExpress                       *bool              `json:"enableExpress"`
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
	t.CountDetails = CountDetailsEntity{}
	t.CountDetails.FromServiceBus(topic.CountDetails)
}

type TopicRequestEntity struct {
	Name    string               `json:"name"`
	Options *TopicRequestOptions `json:"options,omitempty"`
}

func (tr *TopicRequestEntity) GetOptions() (*[]servicebus.TopicManagementOption, *ApiErrorResponse) {
	var opts []servicebus.TopicManagementOption
	opts = make([]servicebus.TopicManagementOption, 0)
	var errorResponse ApiErrorResponse
	if tr.Options != nil {
		if tr.Options.AutoDeleteOnIdle != nil {
			d, err := time.ParseDuration(*tr.Options.AutoDeleteOnIdle)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the AutoDeleteOnIdle from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.TopicWithAutoDeleteOnIdle(&d))
		}
		if tr.Options.EnableBatchedOperation != nil {
			opts = append(opts, servicebus.TopicWithBatchedOperations())
		}
		if tr.Options.EnableDuplicateDetection != nil {
			d, err := time.ParseDuration(*tr.Options.AutoDeleteOnIdle)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the EnableDuplicateDetection from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.TopicWithDuplicateDetection(&d))
		}
		if tr.Options.EnableExpress != nil {
			opts = append(opts, servicebus.TopicWithExpress())
		}
		if tr.Options.MaxSizeInMegabytes != nil {
			opts = append(opts, servicebus.TopicWithMaxSizeInMegabytes(*tr.Options.MaxSizeInMegabytes))
		}
		if tr.Options.DefaultMessageTimeToLive != nil {
			d, err := time.ParseDuration(*tr.Options.AutoDeleteOnIdle)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the DefaultMessageTimeToLive from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.TopicWithMessageTimeToLive(&d))
		}
		if tr.Options.SupportOrdering != nil {
			opts = append(opts, servicebus.TopicWithOrdering())
		}
		if tr.Options.EnablePartitioning != nil {
			opts = append(opts, servicebus.TopicWithPartitioning())
		}
	}

	return &opts, nil
}

func (tr *TopicRequestEntity) IsValidate() (bool, *ApiErrorResponse) {
	var errorResponse ApiErrorResponse

	if tr.Name == "" {
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		return false, &errorResponse
	}

	return true, nil
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
	RequireSession           *bool   `json:"requireSession,omitempty"`
}
