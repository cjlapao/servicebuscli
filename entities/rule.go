package entities

import (
	"bytes"
	"encoding/json"
	"errors"

	servicebus "github.com/Azure/azure-service-bus-go"
)

// RuleRequestEntity struct
type RuleRequestEntity struct {
	Name      string `json:"name"`
	SQLFilter string `json:"sqlFilter"`
	SQLAction string `json:"sqlAction"`
}

// ForwardEntity struct
type ForwardEntity struct {
	To string
	In ForwardingDestinationEntity
}

// ForwardingDestinationEntity Enum
type ForwardingDestinationEntity int

// ForwardingDestination Enum definition
const (
	ForwardToTopic ForwardingDestinationEntity = iota
	ForwardToQueue
)

func (s ForwardingDestinationEntity) String() string {
	return forwardingDestinationToString[s]
}

var forwardingDestinationToString = map[ForwardingDestinationEntity]string{
	ForwardToTopic: "Topic",
	ForwardToQueue: "Queue",
}

var forwardingDestinationToID = map[string]ForwardingDestinationEntity{
	"Topic": ForwardToTopic,
	"topic": ForwardToTopic,
	"Queue": ForwardToQueue,
	"queue": ForwardToQueue,
}

func (s ForwardingDestinationEntity) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(forwardingDestinationToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *ForwardingDestinationEntity) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = forwardingDestinationToID[j]

	return nil
}

type RuleResponseEntity struct {
	Name   string                   `json:"name"`
	ID     string                   `json:"id"`
	Filter RuleResponseFilterEntity `json:"filter"`
	Action RuleResponseActionEntity `json:"action"`
}

func (c *RuleResponseEntity) FromServiceBus(msg *servicebus.RuleEntity) error {
	if msg == nil {
		err := errors.New("entity cannot be null")
		return err
	}

	c.ID = msg.ID
	c.Name = msg.Name
	c.Filter = RuleResponseFilterEntity{}
	c.Action = RuleResponseActionEntity{}

	c.Filter.FromServiceBus(msg.Filter)
	c.Action.FromServiceBus(msg.Action)

	return nil
}

type RuleResponseFilterEntity struct {
	CorrelationID      *string                `json:"correlationID"`
	MessageID          *string                `json:"messageID"`
	To                 *string                `json:"to"`
	ReplyTo            *string                `json:"replyTo"`
	Label              *string                `json:"label"`
	SessionID          *string                `json:"sessionID"`
	ReplyToSessionID   *string                `json:"replyToSessionID"`
	ContentType        *string                `json:"contentType"`
	Properties         map[string]interface{} `json:"properties"`
	Type               string                 `json:"type"`
	SQLExpression      *string                `json:"sqlExpression"`
	CompatibilityLevel int                    `json:"compatibilityLevel"`
}

func (c *RuleResponseFilterEntity) FromServiceBus(filter servicebus.FilterDescription) error {
	c.CorrelationID = filter.CorrelationID
	c.MessageID = filter.MessageID
	c.To = filter.To
	c.ReplyTo = filter.ReplyTo
	c.Label = filter.Label
	c.SessionID = filter.SessionID
	c.ReplyToSessionID = filter.ReplyToSessionID
	c.ContentType = filter.ContentType
	c.Properties = filter.Properties
	c.Type = filter.Type
	c.SQLExpression = filter.SQLExpression
	c.CompatibilityLevel = filter.CompatibilityLevel

	return nil
}

type RuleResponseActionEntity struct {
	Type                  string `json:"type"`
	RequiresPreprocessing bool   `json:"requiresPreprocessing"`
	SQLExpression         string `json:"sqlExpression"`
	CompatibilityLevel    int    `json:"compatibilityLevel"`
}

func (c *RuleResponseActionEntity) FromServiceBus(action *servicebus.ActionDescription) error {
	if action == nil {
		err := errors.New("entity cannot be null")
		return err
	}

	c.Type = action.Type
	c.SQLExpression = action.SQLExpression
	c.CompatibilityLevel = action.CompatibilityLevel

	return nil
}
