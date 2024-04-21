package dto

import "time"

type ErrorResponse struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
}
