package entities

type BulkMessageRequest struct {
	Messages []MessageRequest `json:"messages"`
}

type BulkTemplateMessageRequest struct {
	TotalMessages             int            `json:"totalMessages"`
	Template                  MessageRequest `json:"template"`
	BatchOf                   int            `json:"batchOf"`
	WaitBetweenBatchesInMilli int            `json:"waitBetweenBatchesInMilli"`
}
