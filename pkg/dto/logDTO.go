package dto

import "time"

type LogDTO struct {
	Method    string
	URL       string
	Headers   map[string][]string
	Body      map[string]interface{}
	Status    string
	Timestamp time.Time `json:"@timestamp"`
}
