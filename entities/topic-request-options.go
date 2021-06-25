package entities

type TopicRequestOptions struct {
	AutoDeleteOnIdle         *string `json:"autoDeleteOnIdle,omitempty"`
	EnableBatchedOperation   *bool   `json:"enableBatchedOperation,omitempty"`
	EnableDuplicateDetection *string `json:"enableDuplicateDetection,omitempty"`
	EnableExpress            *bool   `json:"enableExpress,omitempty"`
	MaxSizeInMegabytes       *int    `json:"maxSizeInMegabytes,omitempty"`
	DefaultMessageTimeToLive *string `json:"defaultMessageTimeToLive,omitempty"`
	SupportOrdering          *bool   `json:"supportOrdering,omitempty"`
	EnablePartitioning       *bool   `json:"enablePartitioning,omitempty"`
}
