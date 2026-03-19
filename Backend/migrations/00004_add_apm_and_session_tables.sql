-- +goose Up
-- +goose StatementBegin

-- Create apm_thresholds table for Phase 3: APM Threshold Alerts
CREATE TABLE IF NOT EXISTS apm_thresholds (
  id BIGSERIAL PRIMARY KEY,
  source_id UUID NOT NULL,
  endpoint VARCHAR(255) NOT NULL,
  p95_limit INTEGER NOT NULL DEFAULT 1000,
  p99_limit INTEGER NOT NULL DEFAULT 2000,
  email_notify BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_apm_thresholds_source_id ON apm_thresholds(source_id);
CREATE INDEX IF NOT EXISTS idx_apm_thresholds_endpoint ON apm_thresholds(endpoint);

-- Create sessions table for Phase 2: Session Tracking
CREATE TABLE IF NOT EXISTS sessions (
  id BIGSERIAL PRIMARY KEY,
  user_id VARCHAR(100) NOT NULL,
  source_id UUID NOT NULL,
  auth_method VARCHAR(50) NOT NULL,
  ip_address VARCHAR(45),
  browser VARCHAR(100),
  device VARCHAR(100),
  is_active BOOLEAN NOT NULL DEFAULT true,
  last_activity TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_source_id ON sessions(source_id);
CREATE INDEX IF NOT EXISTS idx_sessions_is_active ON sessions(is_active);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS apm_thresholds;
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
