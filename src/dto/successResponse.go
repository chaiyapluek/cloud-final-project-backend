package dto

import "net/http"

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewSuccessResponse(code int, data any, message ...string) *SuccessResponse {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = http.StatusText(code)
	}
	return &SuccessResponse{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}
