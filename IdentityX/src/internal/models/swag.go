package models

import (
	"time"
)

type ErrorResponse struct {
	Module    string    `json:"module"`
	Message   string    `json:"message"`
	Trace     []string  `json:"trace"`
	Timestamp time.Time `json:"timestamp"`
	Code      int       `json:"code"`
}
