package entities

import (
	"net/http"
	"strings"
	"time"

	azservicebus "github.com/Azure/azure-service-bus-go"
	servicebus "github.com/Azure/azure-service-bus-go"
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

type SubscriptionRequestEntity struct {
	Name              string                      `json:"name"`
	TopicName         string                      `json:"topicName"`
	UserDescription   string                      `json:"userDescription"`
	MaxDeliveryCount  int32                       `json:"maxDeliveryCount,omitempty"`
	Forward           *ForwardEntity              `json:"forward,omitempty"`
	ForwardDeadLetter *ForwardEntity              `json:"forwardDeadLetter,omitempty"`
	Rules             []*RuleRequestEntity        `json:"rules,omitempty"`
	Options           *SubscriptionRequestOptions `json:"options,omitempty"`
}

// NewSubscriptionRequest Creates a new subscription entity
func NewSubscriptionRequest(topicName string, name string) *SubscriptionRequestEntity {
	result := SubscriptionRequestEntity{
		Name:             name,
		TopicName:        topicName,
		MaxDeliveryCount: 10,
	}

	result.Rules = make([]*RuleRequestEntity, 0)
	result.Forward = &ForwardEntity{}
	result.Forward.In = ForwardToTopic
	result.ForwardDeadLetter = &ForwardEntity{}
	result.Options = &SubscriptionRequestOptions{}

	return &result
}

// AddSQLFilter Adds a Sql filter to a specific Rule
func (s *SubscriptionRequestEntity) AddSQLFilter(ruleName string, filter string) {
	var rule RuleRequestEntity
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
		rule = RuleRequestEntity{
			Name:      ruleName,
			SQLFilter: filter,
		}
		s.Rules = append(s.Rules, &rule)
	}
}

// AddSQLAction Adds a Sql Action to a specific rule
func (s *SubscriptionRequestEntity) AddSQLAction(ruleName string, action string) {
	var rule RuleRequestEntity
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
		rule = RuleRequestEntity{
			Name:      ruleName,
			SQLAction: action,
		}
		s.Rules = append(s.Rules, &rule)
	}
}

// MapMessageForwardFlag Maps a forward flag string into it's sub components
func (s *SubscriptionRequestEntity) MapMessageForwardFlag(value string) {
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
func (s *SubscriptionRequestEntity) MapDeadLetterForwardFlag(value string) {
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
func (s *SubscriptionRequestEntity) MapRuleFlag(value string) {
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

func (tr *SubscriptionRequestEntity) GetOptions() (*[]servicebus.SubscriptionManagementOption, *ApiErrorResponse) {
	var opts []servicebus.SubscriptionManagementOption
	opts = make([]servicebus.SubscriptionManagementOption, 0)
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
			opts = append(opts, servicebus.SubscriptionWithAutoDeleteOnIdle(&d))
		}
		if tr.Options.EnableBatchedOperation != nil {
			opts = append(opts, servicebus.SubscriptionWithBatchedOperations())
		}
		if tr.Options.SubscriptionWithDeadLetteringOnMessageExpiration != nil {
			opts = append(opts, servicebus.SubscriptionWithDeadLetteringOnMessageExpiration())
		}
		if tr.Options.SubscriptionWithDeadLetteringOnMessageExpiration != nil {
			opts = append(opts, servicebus.SubscriptionWithDeadLetteringOnMessageExpiration())
		}
		if tr.Options.RequireSession != nil {
			opts = append(opts, servicebus.SubscriptionWithRequiredSessions())
		}
		if tr.Options.DefaultMessageTimeToLive != nil {
			d, err := time.ParseDuration(*tr.Options.DefaultMessageTimeToLive)
			if err != nil {
				errorResponse.Code = http.StatusBadRequest
				errorResponse.Error = "Duration Parse Error"
				errorResponse.Message = "There was an error processing the DefaultMessageTimeToLive from string"
				return nil, &errorResponse
			}
			opts = append(opts, servicebus.SubscriptionWithMessageTimeToLive(&d))
		}
		if tr.Options.LockDuration != nil {
			d, err := time.ParseDuration(*tr.Options.LockDuration)
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

func (tr *SubscriptionRequestEntity) ValidateSubscriptionRequest() (bool, *ApiErrorResponse) {
	var errorResponse ApiErrorResponse

	if tr.Name == "" {
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Subscription name is null"
		errorResponse.Message = "Subscription Name cannot be null"
		return false, &errorResponse
	}

	if tr.TopicName == "" {
		errorResponse.Code = http.StatusBadRequest
		errorResponse.Error = "Topic Name is null"
		errorResponse.Message = "Topic Name cannot be null"
		return false, &errorResponse
	}

	return true, nil
}

type SubscriptionRequestOptions struct {
	AutoDeleteOnIdle                                 *string `json:"autoDeleteOnIdle,omitempty"`
	DefaultMessageTimeToLive                         *string `json:"defaultMessageTimeToLive,omitempty"`
	LockDuration                                     *string `json:"lockDuration,omitempty"`
	EnableBatchedOperation                           *bool   `json:"enableBatchedOperation,omitempty"`
	SubscriptionWithDeadLetteringOnMessageExpiration *bool   `json:"subscriptionWithDeadLetteringOnMessageExpiration,omitempty"`
	RequireSession                                   *bool   `json:"requireSession,omitempty"`
}
