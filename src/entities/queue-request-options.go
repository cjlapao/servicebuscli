package entities

type QueueRequestOptions struct {
	AutoDeleteOnIdle                 *string `json:"autoDeleteOnIdle,omitempty"`
	EnableDuplicateDetection         *string `json:"enableDuplicateDetection,omitempty"`
	MaxSizeInMegabytes               *int    `json:"maxSizeInMegabytes,omitempty"`
	DefaultMessageTimeToLive         *string `json:"defaultMessageTimeToLive,omitempty"`
	LockDuration                     *string `json:"lockDuration,omitempty"`
	SupportOrdering                  *bool   `json:"supportOrdering,omitempty"`
	EnablePartitioning               *bool   `json:"enablePartitioning,omitempty"`
	RequireSession                   *bool   `json:"requireSession,omitempty"`
	DeadLetteringOnMessageExpiration *bool   `json:"deadLetteringOnMessageExpiration,omitempty"`
}
