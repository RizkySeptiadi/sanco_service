package structs

import "time"

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// NewResponse creates a new Response instance
func NewResponse[T any](code int, message string, data T) Response[T] {
	return Response[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

type Sanco_audit struct {
	ID            uint      `gorm:"primarykey"`
	UserID        int64     `json:"user_id"`
	AuditableType string    `json:"auditable_type"`
	AuditableID   uint      `json:"auditable_id"`
	Event         string    `json:"event"`
	OldValues     string    `json:"old_values"`
	NewValues     string    `json:"new_values"`
	URL           string    `json:"url"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	CreatedAt     time.Time `json:"created_at"`
}
