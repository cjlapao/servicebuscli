package entities

// LoginErrorResponse entity
type ApiSuccessResponse struct {
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
