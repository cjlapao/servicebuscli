package entities

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/cjlapao/common-go/helper"
)

type MessageRequest struct {
	Label          string                 `json:"label"`
	CorrelationID  string                 `json:"correlationId"`
	ContentType    string                 `json:"contentType"`
	Data           map[string]interface{} `json:"data"`
	UserProperties map[string]interface{} `json:"userProperties"`
}

func (m *MessageRequest) ToServiceBus() (*servicebus.Message, error) {
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
	} else {
		sbMessage.ContentType = m.ContentType
	}

	if m.CorrelationID != "" {
		sbMessage.CorrelationID = m.CorrelationID
	}

	return &sbMessage, nil
}

func (m *MessageRequest) FromServiceBus(msg *servicebus.Message) error {
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

func (m *MessageRequest) FromFile(filePath string) error {
	fileExists := helper.FileExists(filePath)

	if !fileExists {
		err := errors.New("file " + filePath + " was not found")
		return err
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileContent, m)
	if err != nil {
		return err
	}

	return nil
}
