package contracts

import "time"

type ErrorResponse struct {
	Module    string    `json:"module"`
	Message   string    `json:"message"`
	ErrorID   string    `json:"error_id"`
	Trace     []string  `json:"trace"`
	Timestamp time.Time `json:"timestamp"`
	Code      int       `json:"code"`
}
