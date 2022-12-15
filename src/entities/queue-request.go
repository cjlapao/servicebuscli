package entities

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/common-go/helper"
)

// QueueEntity structure
type QueueRequest struct {
	Name              string               `json:"name"`
	MaxDeliveryCount  int32                `json:"maxDeliveryCount"`
	Forward           *Forward             `json:"forward"`
	ForwardDeadLetter *Forward             `json:"forwardDeadLetter"`
	Options           *QueueRequestOptions `json:"options,omitempty"`
}

// NewQueue Creates a Queue entity
func NewQueueRequest(name string) *QueueRequest {
	result := QueueRequest{
		MaxDeliveryCount: 10,
	}

	result.Name = name
	result.Forward.In = ForwardToQueue
	result.ForwardDeadLetter.In = ForwardToQueue

	return &result
}

// MapMessageForwardFlag Maps a forward flag string into it's sub components
func (q *QueueRequest) MapMessageForwardFlag(value string) {
	if value != "" {
		forwardMapped := strings.Split(value, ":")
		if len(forwardMapped) == 1 {
			q.Forward.To = forwardMapped[0]
		} else if len(forwardMapped) == 2 {
			q.Forward.To = forwardMapped[1]
			switch strings.ToLower(forwardMapped[0]) {
			case "topic":
				q.Forward.In = ForwardToTopic
			case "queue":
				q.Forward.In = ForwardToQueue
			}
		}
	}
}

// MapDeadLetterForwardFlag Maps a forward dead letter flag string into it's sub components
func (q *QueueRequest) MapDeadLetterForwardFlag(value string) {
	if value != "" {
		forwardMapped := strings.Split(value, ":")
		if len(forwardMapped) == 1 {
			q.ForwardDeadLetter.To = forwardMapped[0]
		} else if len(forwardMapped) == 2 {
			q.ForwardDeadLetter.To = forwardMapped[1]
			switch strings.ToLower(forwardMapped[0]) {
			case "topic":
				q.ForwardDeadLetter.In = ForwardToTopic
			case "queue":
				q.ForwardDeadLetter.In = ForwardToQueue
			}
		}
	}
}

func (q *QueueRequest) GetOptions() (*[]servicebus.QueueManagementOption, *ApiErrorResponse) {
	var opts []servicebus.QueueManagementOption
	opts = make([]servicebus.QueueManagementOption, 0)
	var errorResponse ApiErrorResponse

	if q.Options != nil {
		if q.Options.AutoDeleteOnIdle != nil {
			d, err := time.ParseDuration(*q.Options.AutoDeleteOnIdle)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the AutoDeleteOnIdle from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.QueueEntityWithAutoDeleteOnIdle(&d))
		}
		if q.Options.DeadLetteringOnMessageExpiration != nil && *q.Options.DeadLetteringOnMessageExpiration {
			opts = append(opts, servicebus.QueueEntityWithDeadLetteringOnMessageExpiration())
		}
		if q.Options.EnableDuplicateDetection != nil {
			d, err := time.ParseDuration(*q.Options.AutoDeleteOnIdle)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the EnableDuplicateDetection from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.QueueEntityWithDuplicateDetection(&d))
		}
		if q.Options.LockDuration != nil {
			d, err := time.ParseDuration(*q.Options.LockDuration)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the LockDuration from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.QueueEntityWithLockDuration(&d))
		}
		if q.Options.MaxSizeInMegabytes != nil {
			opts = append(opts, servicebus.QueueEntityWithMaxSizeInMegabytes(*q.Options.MaxSizeInMegabytes))
		}
		if q.Options.DefaultMessageTimeToLive != nil {
			d, err := time.ParseDuration(*q.Options.DefaultMessageTimeToLive)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the DefaultMessageTimeToLive from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.QueueEntityWithMessageTimeToLive(&d))
		}
		if q.Options.RequireSession != nil && *q.Options.RequireSession {
			opts = append(opts, servicebus.QueueEntityWithRequiredSessions())
		}
		if q.Options.EnablePartitioning != nil && *q.Options.EnablePartitioning {
			opts = append(opts, servicebus.QueueEntityWithPartitioning())
		}
	}

	return &opts, nil
}

func (q *QueueRequest) IsValid() (bool, *ApiErrorResponse) {
	var errorResponse ApiErrorResponse

	if q.Name == "" {
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Queue name is null"
		errorResponse.Message = "Queue name cannot be null"
		return false, &errorResponse
	}

	_, errResp := q.GetOptions()

	if errResp != nil {
		return false, errResp
	}

	return true, nil
}

func (q *QueueRequest) FromFile(filePath string) error {
	fileExists := helper.FileExists(filePath)

	if !fileExists {
		err := errors.New("file " + filePath + " was not found")
		return err
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileContent, q)
	if err != nil {
		return err
	}

	return nil
}
