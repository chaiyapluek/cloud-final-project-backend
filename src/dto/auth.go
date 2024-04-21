package dto

import "time"

type LoginRequest struct {
	AttemptId string     `json:"attemptId"`
	Email     string     `json:"email"`
	Code      string     `json:"code"`
	RequestAt *time.Time `json:"requestAt"`
}

type RegisterRequest struct {
	AttemptId string     `json:"attemptId"`
	Email     string     `json:"email"`
	Code      string     `json:"code"`
	RequestAt *time.Time `json:"requestAt"`
}

type UserResponse struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
}

type LoginAttemptRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginAttemptResponse struct {
	AttemptId string `json:"attemptId"`
	Code      string `json:"code"`
}

type RegisterAttemptRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RegisterAttemptResponse struct {
	AttemptId string `json:"attemptId"`
	Code      string `json:"code"`
}
