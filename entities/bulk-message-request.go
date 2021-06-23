package entities

type BulkMessageRequest struct {
	Messages []MessageRequest `json:"messages"`
}
