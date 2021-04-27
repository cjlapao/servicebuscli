package entities

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/servicebuscli-go/duration"
)

// QueueResponseEntity
type QueueResponseEntity struct {
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

func (q *QueueResponseEntity) FromServiceBus(queue *servicebus.QueueEntity) {
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

// QueueEntity structure
type QueueRequestEntity struct {
	Name              string               `json:"name"`
	MaxDeliveryCount  int32                `json:"maxDeliveryCount"`
	Forward           *ForwardEntity       `json:"forward"`
	ForwardDeadLetter *ForwardEntity       `json:"forwardDeadLetter"`
	Options           *QueueRequestOptions `json:"options,omitempty"`
}

// NewQueue Creates a Queue entity
func NewQueueRequest(name string) *QueueRequestEntity {
	result := QueueRequestEntity{
		MaxDeliveryCount: 10,
	}

	result.Name = name
	result.Forward.In = ForwardToQueue
	result.ForwardDeadLetter.In = ForwardToQueue

	return &result
}

// MapMessageForwardFlag Maps a forward flag string into it's sub components
func (s *QueueRequestEntity) MapMessageForwardFlag(value string) {
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
func (s *QueueRequestEntity) MapDeadLetterForwardFlag(value string) {
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

func (qr *QueueRequestEntity) GetOptions() (*[]servicebus.QueueManagementOption, *ApiErrorResponse) {
	var opts []servicebus.QueueManagementOption
	opts = make([]servicebus.QueueManagementOption, 0)
	var errorResponse ApiErrorResponse

	if qr.Options != nil {
		if qr.Options.AutoDeleteOnIdle != nil {
			d, err := time.ParseDuration(*qr.Options.AutoDeleteOnIdle)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the AutoDeleteOnIdle from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.QueueEntityWithAutoDeleteOnIdle(&d))
		}
		if qr.Options.DeadLetteringOnMessageExpiration != nil && *qr.Options.DeadLetteringOnMessageExpiration {
			opts = append(opts, servicebus.QueueEntityWithDeadLetteringOnMessageExpiration())
		}
		if qr.Options.EnableDuplicateDetection != nil {
			d, err := time.ParseDuration(*qr.Options.AutoDeleteOnIdle)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the EnableDuplicateDetection from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.QueueEntityWithDuplicateDetection(&d))
		}
		if qr.Options.LockDuration != nil {
			d, err := time.ParseDuration(*qr.Options.AutoDeleteOnIdle)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the LockDuration from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.QueueEntityWithLockDuration(&d))
		}
		if qr.Options.MaxSizeInMegabytes != nil {
			opts = append(opts, servicebus.QueueEntityWithMaxSizeInMegabytes(*qr.Options.MaxSizeInMegabytes))
		}
		if qr.Options.DefaultMessageTimeToLive != nil {
			d, err := time.ParseDuration(*qr.Options.DefaultMessageTimeToLive)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the DefaultMessageTimeToLive from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.QueueEntityWithMessageTimeToLive(&d))
		}
		if qr.Options.RequireSession != nil && *qr.Options.RequireSession {
			opts = append(opts, servicebus.QueueEntityWithRequiredSessions())
		}
		if qr.Options.EnablePartitioning != nil && *qr.Options.EnablePartitioning {
			opts = append(opts, servicebus.QueueEntityWithPartitioning())
		}
	}

	return &opts, nil
}

func (qr *QueueRequestEntity) IsValid() (bool, *ApiErrorResponse) {
	var errorResponse ApiErrorResponse

	if qr.Name == "" {
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Queue name is null"
		errorResponse.Message = "Queue name cannot be null"
		return false, &errorResponse
	}

	_, errResp := qr.GetOptions()

	if errResp != nil {
		return false, errResp
	}

	return true, nil
}

type QueueRequestOptions struct {
	AutoDeleteOnIdle                 *string `json:"autoDeleteOnIdle,omitempty"`
	EnableDuplicateDetection         *string `json:"enableDuplicateDetection,omitempty"`
	MaxSizeInMegabytes               *int    `json:"maxSizeInMegabytes,omitempty"`
	DefaultMessageTimeToLive         *string `json:"defaultMessageTimeToLive,omitempty"`
	LockDuration                     *string `json:"lockDuration,omitempty"`
	SupportOrdering                  *bool   `json:"supportOrdering,omitempty"`
	EnablePartitioning               *bool   `json:"enablePartitioning,omitempty"`
	RequireSession                   *bool   `json:"requireSession,omitempty"`
	DeadLetteringOnMessageExpiration *bool   `json:"deadLetteringOnMessageExpiration,omitempty"`
}
