-- +goose Up
-- Enable pgvector extension (for error embeddings)
CREATE EXTENSION IF NOT EXISTS vector;

-- Create error_embeddings table
CREATE TABLE IF NOT EXISTS error_embeddings (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT UNIQUE,
    fingerprint VARCHAR(255),
    embedding vector(384),
    message_hash VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_error_embeddings_fingerprint ON error_embeddings(fingerprint);
CREATE INDEX IF NOT EXISTS idx_error_embeddings_message_hash ON error_embeddings(message_hash);

-- Create error_clusters table
CREATE TABLE IF NOT EXISTS error_clusters (
    id BIGSERIAL PRIMARY KEY,
    cluster_id VARCHAR(255) UNIQUE,
    representative TEXT,
    message_pattern TEXT,
    count INTEGER DEFAULT 0,
    first_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_error_clusters_cluster_id ON error_clusters(cluster_id);

-- Create cluster_members table
CREATE TABLE IF NOT EXISTS cluster_members (
    id BIGSERIAL PRIMARY KEY,
    cluster_id VARCHAR(255),
    log_id BIGINT,
    distance FLOAT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cluster_members_cluster_id ON cluster_members(cluster_id);
CREATE INDEX IF NOT EXISTS idx_cluster_members_log_id ON cluster_members(log_id);

-- Create config_audit_trails table
CREATE TABLE IF NOT EXISTS config_audit_trails (
    id BIGSERIAL PRIMARY KEY,
    source_id VARCHAR(255),
    environment VARCHAR(255),
    key VARCHAR(255),
    old_value TEXT,
    new_value TEXT,
    changed_by BIGINT,
    change_type VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_config_audit_source_id ON config_audit_trails(source_id);
CREATE INDEX IF NOT EXISTS idx_config_audit_environment ON config_audit_trails(environment);
CREATE INDEX IF NOT EXISTS idx_config_audit_changed_by ON config_audit_trails(changed_by);

-- Create config_webhooks table
CREATE TABLE IF NOT EXISTS config_webhooks (
    id BIGSERIAL PRIMARY KEY,
    source_id VARCHAR(255),
    url TEXT,
    secret TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_config_webhooks_source_id ON config_webhooks(source_id);

-- +goose Down
DROP TABLE IF EXISTS config_webhooks;
DROP TABLE IF EXISTS config_audit_trails;
DROP TABLE IF EXISTS cluster_members;
DROP TABLE IF EXISTS error_clusters;
DROP TABLE IF EXISTS error_embeddings;
DROP EXTENSION IF EXISTS vector;
