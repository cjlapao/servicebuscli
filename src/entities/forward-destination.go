package entities

import (
	"bytes"
	"encoding/json"
)

// ForwardDestination Enum
type ForwardDestination int

// ForwardingDestination Enum definition
const (
	ForwardToTopic ForwardDestination = iota
	ForwardToQueue
)

// String Gets the ForwardingDestinationEntity Enum string representation
func (s ForwardDestination) String() string {
	return forwardingDestinationToString[s]
}

var forwardingDestinationToString = map[ForwardDestination]string{
	ForwardToTopic: "Topic",
	ForwardToQueue: "Queue",
}

var forwardingDestinationToID = map[string]ForwardDestination{
	"Topic": ForwardToTopic,
	"topic": ForwardToTopic,
	"Queue": ForwardToQueue,
	"queue": ForwardToQueue,
}

// MarshalJSON Custom Marsheller of ForwardingDestinationEntity Enum to JSON
func (s ForwardDestination) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(forwardingDestinationToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON Custom UnMarsheller of ForwardingDestinationEntity Enum to JSON
func (s *ForwardDestination) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = forwardingDestinationToID[j]

	return nil
}
