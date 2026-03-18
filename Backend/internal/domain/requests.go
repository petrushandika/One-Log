package domain

// IngestLogRequest presentation request for ingest log
type IngestLogRequest struct {
	// SourceID is derived from API key middleware (clients should not send it).
	// Kept optional for backward compatibility with older clients/tests.
	SourceID string `json:"source_id" binding:"omitempty,max=50"`

	// Categories aligned with Documentation (MVP + roadmap).
	Category string `json:"category" binding:"required,oneof=AUTH_EVENT USER_ACTIVITY SYSTEM_ERROR SECURITY PERFORMANCE AUDIT_TRAIL"`

	// Levels aligned with Documentation. Accept WARNING as legacy alias for WARN.
	Level      string                 `json:"level" binding:"required,oneof=DEBUG INFO WARN WARNING ERROR CRITICAL"`
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
	// Optional: health check endpoint for uptime worker.
	HealthURL string `json:"health_url" binding:"omitempty,url,max=255"`
}

// UpdateSourceRequest allows updating mutable source fields.
type UpdateSourceRequest struct {
	Name      *string `json:"name" binding:"omitempty,min=3,max=100"`
	HealthURL *string `json:"health_url" binding:"omitempty,url,max=255"`
	Status    *string `json:"status" binding:"omitempty,oneof=ONLINE OFFLINE DEGRADED MAINTENANCE"`
}
