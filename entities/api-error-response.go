package entities

// LoginErrorResponse entity
type ApiErrorResponse struct {
	Code    int32  `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
