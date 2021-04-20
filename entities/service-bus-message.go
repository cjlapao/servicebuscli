package entities

import (
	"encoding/json"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type ServiceBusMessage struct {
	Label          string                 `json:"label"`
	CorrelationID  string                 `json:"correlationId"`
	ContentType    string                 `json:"contentType"`
	Data           map[string]interface{} `json:"data"`
	UserProperties map[string]interface{} `json:"userProperties"`
}

func (m *ServiceBusMessage) ToServiceBus() (*servicebus.Message, error) {
	messageData, err := json.MarshalIndent(m.Data, "", "  ")
	if err != nil {
		return nil, err
	}

	sbMessage := servicebus.Message{
		Data:           messageData,
		UserProperties: m.UserProperties,
	}

	if m.Label != "" {
		sbMessage.Label = m.Label
	}

	if m.ContentType == "" {
		sbMessage.ContentType = "application/json"
	}

	if m.CorrelationID != "" {
		sbMessage.CorrelationID = m.CorrelationID
	}

	return &sbMessage, nil
}
