package domain

type IngestLogRequest struct {
	SourceID   string                 `json:"source_id" binding:"required,max=50"`
	Category   string                 `json:"category" binding:"required,oneof=AUTH_EVENT USER_ACTIVITY SYSTEM_ERROR"`
	Level      string                 `json:"level" binding:"required,oneof=INFO WARNING ERROR CRITICAL"`
	Message    string                 `json:"message" binding:"required"`
	StackTrace string                 `json:"stack_trace"`
	Context    map[string]interface{} `json:"context"`
	IPAddress  string                 `json:"ip_address" binding:"omitempty,ip"`
}
