-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  name VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS sources (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id BIGINT NOT NULL,
  name VARCHAR(100) NOT NULL,
  api_key VARCHAR(255) NOT NULL UNIQUE,
  health_url VARCHAR(255),
  status VARCHAR(20) NOT NULL DEFAULT 'ONLINE',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_sources_user_id ON sources(user_id);

CREATE TABLE IF NOT EXISTS log_entries (
  id BIGSERIAL PRIMARY KEY,
  source_id UUID NOT NULL,
  category VARCHAR(50) NOT NULL,
  level VARCHAR(20) NOT NULL,
  message TEXT NOT NULL,
  context JSONB,
  stack_trace TEXT,
  ip_address VARCHAR(45),
  ai_insight JSONB,
  fingerprint VARCHAR(64),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_log_entries_source_id ON log_entries(source_id);
CREATE INDEX IF NOT EXISTS idx_log_entries_category ON log_entries(category);
CREATE INDEX IF NOT EXISTS idx_log_entries_level ON log_entries(level);
CREATE INDEX IF NOT EXISTS idx_log_entries_created_at ON log_entries(created_at);
CREATE INDEX IF NOT EXISTS idx_log_entries_fingerprint ON log_entries(fingerprint);
CREATE INDEX IF NOT EXISTS idx_log_entries_context_gin ON log_entries USING GIN (context);

CREATE TABLE IF NOT EXISTS issues (
  fingerprint VARCHAR(64) PRIMARY KEY,
  source_id UUID NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
  category VARCHAR(50),
  level VARCHAR(20),
  message_sample TEXT,
  occurrence_count BIGINT NOT NULL DEFAULT 1,
  first_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_issues_source_id ON issues(source_id);
CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
CREATE INDEX IF NOT EXISTS idx_issues_last_seen_at ON issues(last_seen_at);

CREATE TABLE IF NOT EXISTS source_configs (
  id BIGSERIAL PRIMARY KEY,
  source_id UUID NOT NULL,
  environment VARCHAR(30) NOT NULL DEFAULT 'production',
  key VARCHAR(100) NOT NULL,
  value TEXT,
  is_secret BOOLEAN NOT NULL DEFAULT FALSE,
  updated_by BIGINT,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_source_configs_source_id ON source_configs(source_id);
CREATE INDEX IF NOT EXISTS idx_source_configs_env_key ON source_configs(environment, key);
CREATE UNIQUE INDEX IF NOT EXISTS uidx_source_configs_source_env_key ON source_configs(source_id, environment, key);

CREATE TABLE IF NOT EXISTS source_config_histories (
  id BIGSERIAL PRIMARY KEY,
  source_id UUID NOT NULL,
  environment VARCHAR(30) NOT NULL,
  key VARCHAR(100) NOT NULL,
  value TEXT,
  is_secret BOOLEAN NOT NULL DEFAULT FALSE,
  version BIGINT NOT NULL DEFAULT 1,
  updated_by BIGINT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sch_source_env_key ON source_config_histories(source_id, environment, key);
CREATE INDEX IF NOT EXISTS idx_sch_version ON source_config_histories(version);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS source_config_histories;
DROP TABLE IF EXISTS source_configs;
DROP TABLE IF EXISTS issues;
DROP TABLE IF EXISTS log_entries;
DROP TABLE IF EXISTS sources;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

