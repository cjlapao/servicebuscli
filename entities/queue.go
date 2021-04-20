package entities

import (
	"fmt"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/servicebuscli-go/duration"
)

// QueueEntity
type QueueEntity struct {
	Name                                string             `json:"name"`
	ID                                  string             `json:"id"`
	CountDetails                        CountDetailsEntity `json:"countDetails"`
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

func (q *QueueEntity) FromServiceBus(queue *servicebus.QueueEntity) {
	if queue == nil {
		return
	}

	fmt.Println(*queue.AutoDeleteOnIdle)
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
	q.CountDetails = CountDetailsEntity{}
	q.CountDetails.FromServiceBus(queue.CountDetails)
}

// type QueueRequestEntity struct {
// 	Name    string               `json:"name"`
// 	Options *QueueRequestOptions `json:"options,omitempty"`
// }

// func (qr *QueueRequestEntity) GetOptions() (*[]servicebus.QueueManagementOption, *ApiErrorResponse) {
// 	var opts []servicebus.QueueManagementOption
// 	opts = make([]servicebus.QueueManagementOption, 0)
// 	var errorResponse ApiErrorResponse

// 	if qr.Options != nil {
// 		if qr.Options.ForwardTo != nil {
// 			opts = append(opts, servicebus.qu)
// 		}
// 		if qr.Options.AutoDeleteOnIdle != nil {
// 			d, err := time.ParseDuration(*qr.Options.AutoDeleteOnIdle)
// 			if err != nil {
// 				errorResponse.Code = http.StatusBadRequest
// 				errorResponse.Error = "Duration Parse Error"
// 				errorResponse.Message = "There was an error processing the AutoDeleteOnIdle from string"
// 				return nil, &errorResponse
// 			}
// 			opts = append(opts, servicebus.TopicWithAutoDeleteOnIdle(&d))
// 		}
// 		if qr.Options.EnableBatchedOperation != nil {
// 			opts = append(opts, servicebus.TopicWithBatchedOperations())
// 		}
// 		if qr.Options.EnableDuplicateDetection != nil {
// 			d, err := time.ParseDuration(*qr.Options.AutoDeleteOnIdle)
// 			if err != nil {
// 				errorResponse.Code = http.StatusBadRequest
// 				errorResponse.Error = "Duration Parse Error"
// 				errorResponse.Message = "There was an error processing the EnableDuplicateDetection from string"
// 				return nil, &errorResponse
// 			}
// 			opts = append(opts, servicebus.TopicWithDuplicateDetection(&d))
// 		}
// 		if qr.Options.EnableExpress != nil {
// 			opts = append(opts, servicebus.TopicWithExpress())
// 		}
// 		if qr.Options.MaxSizeInMegabytes != nil {
// 			opts = append(opts, servicebus.TopicWithMaxSizeInMegabytes(*qr.Options.MaxSizeInMegabytes))
// 		}
// 		if qr.Options.DefaultMessageTimeToLive != nil {
// 			d, err := time.ParseDuration(*qr.Options.AutoDeleteOnIdle)
// 			if err != nil {
// 				errorResponse.Code = http.StatusBadRequest
// 				errorResponse.Error = "Duration Parse Error"
// 				errorResponse.Message = "There was an error processing the DefaultMessageTimeToLive from string"
// 				return nil, &errorResponse
// 			}
// 			opts = append(opts, servicebus.TopicWithMessageTimeToLive(&d))
// 		}
// 		if qr.Options.SupportOrdering != nil {
// 			opts = append(opts, servicebus.TopicWithOrdering())
// 		}
// 		if qr.Options.EnablePartitioning != nil {
// 			opts = append(opts, servicebus.TopicWithPartitioning())
// 		}
// 	}

// 	return &opts, nil
// }

// func (tr *QueueRequestEntity) IsValidate() (bool, *ApiErrorResponse) {
// 	var errorResponse ApiErrorResponse

// 	if tr.Name == "" {
// 		errorResponse.Code = http.StatusBadRequest
// 		errorResponse.Error = "Topic name is null"
// 		errorResponse.Message = "Topic name cannot be null"
// 		return false, &errorResponse
// 	}

// 	return true, nil
// }

// type QueueRequestOptions struct {
// 	ForwardTo                     *string `json:"forwardTo"`
// 	ForwardDeadLetteredMessagesTo *string `json:"forwardDeadLetteredMessagesTo"`
// 	AutoDeleteOnIdle              *string `json:"autoDeleteOnIdle,omitempty"`
// 	EnableBatchedOperation        *bool   `json:"enableBatchedOperation,omitempty"`
// 	EnableDuplicateDetection      *bool   `json:"enableDuplicateDetection,omitempty"`
// 	EnableExpress                 *bool   `json:"enableExpress,omitempty"`
// 	MaxSizeInMegabytes            *int    `json:"maxSizeInMegabytes,omitempty"`
// 	DefaultMessageTimeToLive      *string `json:"defaultMessageTimeToLive,omitempty"`
// 	SupportOrdering               *bool   `json:"supportOrdering,omitempty"`
// 	EnablePartitioning            *bool   `json:"enablePartitioning,omitempty"`
// }