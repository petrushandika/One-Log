-- +goose Up
-- +goose StatementBegin

-- 1. Add user_id to sources (nullable first, populate, then enforce NOT NULL)
ALTER TABLE sources ADD COLUMN IF NOT EXISTS user_id BIGINT;
UPDATE sources SET user_id = 1 WHERE user_id IS NULL;
ALTER TABLE sources ALTER COLUMN user_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_sources_user_id ON sources(user_id);

-- 2. Create issues table (Phase 5 — error grouping)
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

-- 3. Create source_config_histories table (Phase 6 — config versioning)
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

-- 4. Create source_configs table if not exists (Phase 6)
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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS source_config_histories;
DROP TABLE IF EXISTS source_configs;
DROP TABLE IF EXISTS issues;
ALTER TABLE sources DROP COLUMN IF EXISTS user_id;
-- +goose StatementEnd
