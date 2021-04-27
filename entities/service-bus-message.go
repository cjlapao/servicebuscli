package entities

import (
	"encoding/json"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type ServiceBusMessageRequest struct {
	Label          string                 `json:"label"`
	CorrelationID  string                 `json:"correlationId"`
	ContentType    string                 `json:"contentType"`
	Data           map[string]interface{} `json:"data"`
	UserProperties map[string]interface{} `json:"userProperties"`
}

func (m *ServiceBusMessageRequest) ToServiceBus() (*servicebus.Message, error) {
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

func (m *ServiceBusMessageRequest) FromServiceBus(msg *servicebus.Message) error {
	m.Data = map[string]interface{}{}
	err := json.Unmarshal(msg.Data, &m.Data)
	if err != nil {
		return err
	}
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
