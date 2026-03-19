package service

import (
	"crypto/sha256"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/petrushandika/one-log/internal/domain"
	"gorm.io/gorm"
)

// ErrorDeduplicationService handles semantic error grouping using embeddings
type ErrorDeduplicationService struct {
	db *gorm.DB
}

func NewErrorDeduplicationService(db *gorm.DB) *ErrorDeduplicationService {
	return &ErrorDeduplicationService{db: db}
}

// GenerateEmbedding creates a simple vector representation of error message
// In production, this should call an external embedding API (OpenAI, etc)
func (s *ErrorDeduplicationService) GenerateEmbedding(message string) []float32 {
	// Simplified: Create a bag-of-words vector
	// In production, use: OpenAI, Cohere, or local embedding model
	words := s.tokenize(message)
	vector := make([]float32, 384)

	for i, word := range words {
		if i >= 384 {
			break
		}
		// Simple hash-based embedding
		hash := sha256.Sum256([]byte(word))
		vector[i] = float32(hash[0]) / 255.0
	}

	// Normalize vector
	s.normalize(vector)

	return vector
}

// FindSimilarErrors finds errors similar to the given embedding
func (s *ErrorDeduplicationService) FindSimilarErrors(embedding []float32, threshold float64) ([]domain.ErrorEmbedding, error) {
	var embeddings []domain.ErrorEmbedding

	// Get all embeddings (in production, use vector DB like pgvector with similarity search)
	if err := s.db.Find(&embeddings).Error; err != nil {
		return nil, err
	}

	var similar []domain.ErrorEmbedding
	for _, emb := range embeddings {
		similarity := s.cosineSimilarity(embedding, emb.Embedding)
		if similarity >= threshold {
			similar = append(similar, emb)
		}
	}

	return similar, nil
}

// CreateErrorCluster creates a new cluster for similar errors
func (s *ErrorDeduplicationService) CreateErrorCluster(messages []string) (*domain.ErrorCluster, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	// Generate pattern from messages
	pattern := s.extractPattern(messages)

	cluster := &domain.ErrorCluster{
		ClusterID:      s.generateClusterID(),
		Representative: messages[0],
		MessagePattern: pattern,
		Count:          len(messages),
	}

	if err := s.db.Create(cluster).Error; err != nil {
		return nil, err
	}

	return cluster, nil
}

// GetErrorClusters returns all error clusters
func (s *ErrorDeduplicationService) GetErrorClusters(limit int) ([]domain.ErrorCluster, error) {
	var clusters []domain.ErrorCluster
	err := s.db.Order("count DESC").Limit(limit).Find(&clusters).Error
	return clusters, err
}

// Helper functions

func (s *ErrorDeduplicationService) tokenize(message string) []string {
	// Simple tokenization
	message = strings.ToLower(message)
	message = strings.ReplaceAll(message, ",", " ")
	message = strings.ReplaceAll(message, ".", " ")
	message = strings.ReplaceAll(message, "!", " ")
	message = strings.ReplaceAll(message, "?", " ")
	return strings.Fields(message)
}

func (s *ErrorDeduplicationService) normalize(vector []float32) {
	var sum float64
	for _, v := range vector {
		sum += float64(v * v)
	}
	norm := math.Sqrt(sum)
	if norm > 0 {
		for i := range vector {
			vector[i] = float32(float64(vector[i]) / norm)
		}
	}
}

func (s *ErrorDeduplicationService) cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float64
	var normA float64
	var normB float64

	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (s *ErrorDeduplicationService) extractPattern(messages []string) string {
	// Extract common pattern from messages
	// Replace dynamic parts with wildcards
	if len(messages) == 0 {
		return ""
	}

	pattern := messages[0]

	// Replace common dynamic values
	replacements := []string{
		"[0-9]+",     // Numbers
		"[a-zA-Z]+",  // Words
		"[0-9a-f-]+", // UUIDs
	}

	for _, repl := range replacements {
		pattern = strings.ReplaceAll(pattern, repl, "*")
	}

	return pattern
}

func (s *ErrorDeduplicationService) generateClusterID() string {
	// Generate unique cluster ID
	return fmt.Sprintf("cluster_%d", time.Now().UnixNano())
}
