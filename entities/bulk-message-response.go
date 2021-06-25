package entities

type BulkMessageResponse struct {
	SuccessCount int `json:"successCount"`
	ErrorCount   int `json:"errorCount"`
}
