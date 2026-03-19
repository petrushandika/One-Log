package domain

import "time"

// ErrorEmbedding stores vector embeddings for error messages
// Note: This table requires pgvector extension
// If pgvector is not available, this table will not be created
type ErrorEmbedding struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	LogID       uint      `json:"log_id" gorm:"uniqueIndex"`
	Fingerprint string    `json:"fingerprint" gorm:"index"`
	Embedding   []float32 `json:"embedding" gorm:"type:float[]"` // Fallback to float array if vector not available
	MessageHash string    `json:"message_hash" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
}

// ErrorCluster represents a group of semantically similar errors
type ErrorCluster struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	ClusterID      string    `json:"cluster_id" gorm:"uniqueIndex"`
	Representative string    `json:"representative"`  // Sample message
	MessagePattern string    `json:"message_pattern"` // Regex-like pattern
	Count          int       `json:"count"`
	FirstSeenAt    time.Time `json:"first_seen_at"`
	LastSeenAt     time.Time `json:"last_seen_at"`
}

// ClusterMember links errors to clusters
type ClusterMember struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ClusterID string    `json:"cluster_id" gorm:"index"`
	LogID     uint      `json:"log_id" gorm:"index"`
	Distance  float64   `json:"distance"` // Cosine distance
	CreatedAt time.Time `json:"created_at"`
}
