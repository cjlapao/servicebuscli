package entities

// LoginErrorResponse entity
type ApiErrorResponse struct {
	Code    int64  `json:"code,omitempty"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// NewApiSuccessResponse Creates a new API Success Response struct
func NewApiErrorResponse(code int64, err string, message string) *ApiErrorResponse {
	if err == "" {
		err = "Server Error"
	}
	result := ApiErrorResponse{
		Code:    code,
		Error:   err,
		Message: message,
	}

	return &result
}
