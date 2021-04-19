package entities

import servicebus "github.com/Azure/azure-service-bus-go"

type CountDetailsEntity struct {
	ActiveMessageCount             *int32 `json:"activeMessageCount"`
	DeadLetterMessageCount         *int32 `json:"deadLetterMessageCount"`
	ScheduledMessageCount          *int32 `json:"scheduledMessageCount"`
	TransferDeadLetterMessageCount *int32 `json:"transferDeadLetterMessageCount"`
	TransferMessageCount           *int32 `json:"transferMessageCount"`
}

func (c *CountDetailsEntity) FromServiceBus(countDetails *servicebus.CountDetails) {
	if countDetails == nil {
		return
	}
	c.ActiveMessageCount = countDetails.ActiveMessageCount
	c.DeadLetterMessageCount = countDetails.DeadLetterMessageCount
	c.ScheduledMessageCount = countDetails.ScheduledMessageCount
	c.TransferDeadLetterMessageCount = countDetails.TransferDeadLetterMessageCount
	c.TransferMessageCount = countDetails.TransferMessageCount
}
