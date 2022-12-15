package entities

import (
	"errors"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type RuleResponseAction struct {
	Type                  string `json:"type"`
	RequiresPreprocessing bool   `json:"requiresPreprocessing"`
	SQLExpression         string `json:"sqlExpression"`
	CompatibilityLevel    int    `json:"compatibilityLevel"`
}

func (e *RuleResponseAction) FromServiceBus(action *servicebus.ActionDescription) error {
	if action == nil {
		err := errors.New("entity cannot be null")
		return err
	}

	e.Type = action.Type
	e.SQLExpression = action.SQLExpression
	e.CompatibilityLevel = action.CompatibilityLevel

	return nil
}
