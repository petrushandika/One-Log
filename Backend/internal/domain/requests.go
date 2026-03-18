package domain

// IngestLogRequest presentation request for ingest log
type IngestLogRequest struct {
	SourceID   string                 `json:"source_id" binding:"required,max=50"`
	Category   string                 `json:"category" binding:"required,oneof=AUTH_EVENT USER_ACTIVITY SYSTEM_ERROR"`
	Level      string                 `json:"level" binding:"required,oneof=INFO WARNING ERROR CRITICAL"`
	Message    string                 `json:"message" binding:"required"`
	StackTrace string                 `json:"stack_trace" binding:"omitempty"`
	Context    map[string]interface{} `json:"context"`
	IPAddress  string                 `json:"ip_address" binding:"omitempty,ip"`
}

// LoginRequest represents the JSON body for admin login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateSourceRequest represents the JSON body to create a log sender app
type CreateSourceRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}
