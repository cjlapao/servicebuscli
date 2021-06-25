package entities

import (
	"encoding/json"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type MessageResponse struct {
	ID             string                 `json:"id"`
	Label          string                 `json:"label"`
	CorrelationID  string                 `json:"correlationId"`
	ContentType    string                 `json:"contentType"`
	Data           map[string]interface{} `json:"data"`
	UserProperties map[string]interface{} `json:"userProperties"`
}

func (m *MessageResponse) FromServiceBus(msg *servicebus.Message) error {
	m.Data = map[string]interface{}{}
	err := json.Unmarshal(msg.Data, &m.Data)
	if err != nil {
		return err
	}

	m.ID = msg.ID
	m.UserProperties = msg.UserProperties

	if msg.Label != "" {
		m.Label = msg.Label
	}

	if msg.ContentType != "" {
		m.ContentType = msg.ContentType
	}

	if msg.CorrelationID != "" {
		m.CorrelationID = msg.CorrelationID
	}

	return nil
}
