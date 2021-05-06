package entities

// ApiSuccessResponse entity
type ApiSuccessResponse struct {
	Code    int64                  `json:"code,omitempty"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// NewApiSuccessResponse Creates a new API Success Response struct
func NewApiSuccessResponse(code int64, message string, data map[string]interface{}) *ApiSuccessResponse {
	result := ApiSuccessResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}

	return &result
}
