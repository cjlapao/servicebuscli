package entities

import (
	"errors"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type RuleResponse struct {
	Name   string             `json:"name"`
	ID     string             `json:"id"`
	Filter RuleResponseFilter `json:"filter"`
	Action RuleResponseAction `json:"action"`
}

func (r *RuleResponse) FromServiceBus(msg *servicebus.RuleEntity) error {
	if msg == nil {
		err := errors.New("entity cannot be null")
		return err
	}

	r.ID = msg.ID
	r.Name = msg.Name
	r.Filter = RuleResponseFilter{}
	r.Action = RuleResponseAction{}

	r.Filter.FromServiceBus(msg.Filter)
	r.Action.FromServiceBus(msg.Action)

	return nil
}
