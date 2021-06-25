package entities

import servicebus "github.com/Azure/azure-service-bus-go"

type RuleResponseFilter struct {
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

func (c *RuleResponseFilter) FromServiceBus(filter servicebus.FilterDescription) error {
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
