package entities

import (
	"net/http"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

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
		if tr.Options.EnableBatchedOperation != nil && *tr.Options.EnableBatchedOperation {
			opts = append(opts, servicebus.TopicWithBatchedOperations())
		}
		if tr.Options.EnableDuplicateDetection != nil {
			d, err := time.ParseDuration(*tr.Options.EnableDuplicateDetection)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the EnableDuplicateDetection from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.TopicWithDuplicateDetection(&d))
		}
		if tr.Options.EnableExpress != nil && *tr.Options.EnableExpress {
			opts = append(opts, servicebus.TopicWithExpress())
		}
		if tr.Options.MaxSizeInMegabytes != nil {
			opts = append(opts, servicebus.TopicWithMaxSizeInMegabytes(*tr.Options.MaxSizeInMegabytes))
		}
		if tr.Options.DefaultMessageTimeToLive != nil {
			d, err := time.ParseDuration(*tr.Options.DefaultMessageTimeToLive)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the DefaultMessageTimeToLive from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.TopicWithMessageTimeToLive(&d))
		}
		if tr.Options.SupportOrdering != nil && *tr.Options.SupportOrdering {
			opts = append(opts, servicebus.TopicWithOrdering())
		}
		if tr.Options.EnablePartitioning != nil && *tr.Options.EnablePartitioning {
			opts = append(opts, servicebus.TopicWithPartitioning())
		}
	}

	return &opts, nil
}

func (tr *TopicRequestEntity) IsValid() (bool, *ApiErrorResponse) {
	var errorResponse ApiErrorResponse

	if tr.Name == "" {
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic name is null"
		errorResponse.Message = "Topic name cannot be null"
		return false, &errorResponse
	}

	return true, nil
}
