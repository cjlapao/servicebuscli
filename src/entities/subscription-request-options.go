package entities

// SubscriptionRequestOptions Entity
type SubscriptionRequestOptions struct {
	AutoDeleteOnIdle                 *string `json:"autoDeleteOnIdle,omitempty"`
	DefaultMessageTimeToLive         *string `json:"defaultMessageTimeToLive,omitempty"`
	LockDuration                     *string `json:"lockDuration,omitempty"`
	EnableBatchedOperation           *bool   `json:"enableBatchedOperation,omitempty"`
	DeadLetteringOnMessageExpiration *bool   `json:"deadLetteringOnMessageExpiration,omitempty"`
	RequireSession                   *bool   `json:"requireSession,omitempty"`
}
