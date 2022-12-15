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

type SubscriptionRequest struct {
	Name              string                      `json:"name"`
	TopicName         string                      `json:"topicName"`
	UserDescription   string                      `json:"userDescription"`
	MaxDeliveryCount  int32                       `json:"maxDeliveryCount,omitempty"`
	Forward           *Forward                    `json:"forward,omitempty"`
	ForwardDeadLetter *Forward                    `json:"forwardDeadLetter,omitempty"`
	Rules             []*RuleRequest              `json:"rules,omitempty"`
	Options           *SubscriptionRequestOptions `json:"options,omitempty"`
}

// NewSubscriptionRequest Creates a new subscription entity
func NewSubscriptionRequest(topicName string, name string) *SubscriptionRequest {
	result := SubscriptionRequest{
		Name:             name,
		TopicName:        topicName,
		MaxDeliveryCount: 10,
	}

	result.Rules = make([]*RuleRequest, 0)
	result.Forward = &Forward{}
	result.Forward.In = ForwardToTopic
	result.ForwardDeadLetter = &Forward{}
	result.Options = &SubscriptionRequestOptions{}

	return &result
}

// AddSQLFilter Adds a Sql filter to a specific Rule
func (s *SubscriptionRequest) AddSQLFilter(ruleName string, filter string) {
	var rule RuleRequest
	ruleFound := false

	for i := range s.Rules {
		if s.Rules[i].Name == ruleName {
			ruleFound = true
			if len(s.Rules[i].SQLFilter) > 0 {
				s.Rules[i].SQLFilter += " "
			}
			s.Rules[i].SQLFilter += filter
			break
		}
	}

	if !ruleFound {
		rule = RuleRequest{
			Name:      ruleName,
			SQLFilter: filter,
		}
		s.Rules = append(s.Rules, &rule)
	}
}

// AddSQLAction Adds a Sql Action to a specific rule
func (s *SubscriptionRequest) AddSQLAction(ruleName string, action string) {
	var rule RuleRequest
	ruleFound := false

	for i := range s.Rules {
		if s.Rules[i].Name == ruleName {
			ruleFound = true
			if len(s.Rules[i].SQLAction) > 0 {
				s.Rules[i].SQLAction += " "
			}
			s.Rules[i].SQLAction += action
			break
		}
	}

	if !ruleFound {
		rule = RuleRequest{
			Name:      ruleName,
			SQLAction: action,
		}
		s.Rules = append(s.Rules, &rule)
	}
}

// MapMessageForwardFlag Maps a forward flag string into it's sub components
func (s *SubscriptionRequest) MapMessageForwardFlag(value string) {
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
func (s *SubscriptionRequest) MapDeadLetterForwardFlag(value string) {
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

// MapRuleFlag Maps a rule flag string into it's sub components
func (s *SubscriptionRequest) MapRuleFlag(value string) {
	if value != "" {
		ruleMapped := strings.Split(value, ":")
		if len(ruleMapped) > 1 {
			s.AddSQLFilter(ruleMapped[0], ruleMapped[1])
			if len(ruleMapped) == 3 {
				s.AddSQLAction(ruleMapped[0], ruleMapped[2])
			}
		}
	}
}

func (s *SubscriptionRequest) GetOptions() (*[]servicebus.SubscriptionManagementOption, *ApiErrorResponse) {
	var opts []servicebus.SubscriptionManagementOption
	opts = make([]servicebus.SubscriptionManagementOption, 0)
	var errorResponse ApiErrorResponse
	if s.Options != nil {
		if s.Options.AutoDeleteOnIdle != nil {
			d, err := time.ParseDuration(*s.Options.AutoDeleteOnIdle)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the AutoDeleteOnIdle from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.SubscriptionWithAutoDeleteOnIdle(&d))
		}
		if s.Options.EnableBatchedOperation != nil && *s.Options.EnableBatchedOperation {
			opts = append(opts, servicebus.SubscriptionWithBatchedOperations())
		}
		if s.Options.DeadLetteringOnMessageExpiration != nil && *s.Options.DeadLetteringOnMessageExpiration {
			opts = append(opts, servicebus.SubscriptionWithDeadLetteringOnMessageExpiration())
		}
		if s.Options.RequireSession != nil && *s.Options.RequireSession {
			opts = append(opts, servicebus.SubscriptionWithRequiredSessions())
		}
		if s.Options.DefaultMessageTimeToLive != nil {
			d, err := time.ParseDuration(*s.Options.DefaultMessageTimeToLive)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the DefaultMessageTimeToLive from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.SubscriptionWithMessageTimeToLive(&d))
		}
		if s.Options.LockDuration != nil {
			d, err := time.ParseDuration(*s.Options.LockDuration)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the LockDuration from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.SubscriptionWithLockDuration(&d))
		}
	}

	return &opts, nil
}

func (s *SubscriptionRequest) IsValid() (bool, *ApiErrorResponse) {
	var errorResponse ApiErrorResponse

	if s.Name == "" {
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Subscription name is null"
		errorResponse.Message = "Subscription Name cannot be null"
		return false, &errorResponse
	}

	if s.TopicName == "" {
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic Name is null"
		errorResponse.Message = "Topic Name cannot be null"
		return false, &errorResponse
	}

	_, errResp := s.GetOptions()

	if errResp != nil {
		return false, errResp
	}

	return true, nil
}

func (s *SubscriptionRequest) FromFile(filePath string) error {
	fileExists := helper.FileExists(filePath)

	if !fileExists {
		err := errors.New("file " + filePath + " was not found")
		return err
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileContent, s)
	if err != nil {
		return err
	}

	return nil
}
