-- +goose Up
-- Add error_embeddings and error_clusters tables for semantic error search
-- Note: pgvector extension is optional; falls back to float[] if not available.

CREATE TABLE IF NOT EXISTS error_embeddings (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT NOT NULL UNIQUE,
    fingerprint VARCHAR(64),
    embedding FLOAT[],
    message_hash VARCHAR(64),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_error_embeddings_fingerprint ON error_embeddings(fingerprint);
CREATE INDEX IF NOT EXISTS idx_error_embeddings_message_hash ON error_embeddings(message_hash);
CREATE INDEX IF NOT EXISTS idx_error_embeddings_log_id ON error_embeddings(log_id);

CREATE TABLE IF NOT EXISTS error_clusters (
    id BIGSERIAL PRIMARY KEY,
    cluster_id VARCHAR(255) NOT NULL UNIQUE,
    representative TEXT,
    message_pattern TEXT,
    count INTEGER NOT NULL DEFAULT 0,
    first_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_error_clusters_cluster_id ON error_clusters(cluster_id);

CREATE TABLE IF NOT EXISTS cluster_members (
    id BIGSERIAL PRIMARY KEY,
    cluster_id VARCHAR(255) NOT NULL,
    log_id BIGINT NOT NULL,
    distance DOUBLE PRECISION,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cluster_members_cluster_id ON cluster_members(cluster_id);
CREATE INDEX IF NOT EXISTS idx_cluster_members_log_id ON cluster_members(log_id);

-- +goose Down
DROP TABLE IF EXISTS cluster_members;
DROP TABLE IF EXISTS error_clusters;
DROP TABLE IF EXISTS error_embeddings;
