-- +goose Up
-- Complete migration to ensure all tables exist

-- Fix issues table - add missing columns if not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='issues' AND column_name='resolved_at') THEN
        ALTER TABLE issues ADD COLUMN resolved_at TIMESTAMP WITH TIME ZONE;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='issues' AND column_name='is_regression') THEN
        ALTER TABLE issues ADD COLUMN is_regression BOOLEAN DEFAULT FALSE;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='issues' AND column_name='regression_alert_sent') THEN
        ALTER TABLE issues ADD COLUMN regression_alert_sent BOOLEAN DEFAULT FALSE;
    END IF;
END $$;

-- Create sessions table if not exists
CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    source_id UUID NOT NULL,
    auth_method VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45),
    browser VARCHAR(100),
    device VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_source_id ON sessions(source_id);
CREATE INDEX IF NOT EXISTS idx_sessions_is_active ON sessions(is_active);

-- Create activity_feeds table if not exists
CREATE TABLE IF NOT EXISTS activity_feeds (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    source_id VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(255),
    context JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_activity_feeds_user_id ON activity_feeds(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_feeds_source_id ON activity_feeds(source_id);
CREATE INDEX IF NOT EXISTS idx_activity_feeds_action ON activity_feeds(action);
CREATE INDEX IF NOT EXISTS idx_activity_feeds_created_at ON activity_feeds(created_at DESC);

-- Create compliance_exports table if not exists
CREATE TABLE IF NOT EXISTS compliance_exports (
    id BIGSERIAL PRIMARY KEY,
    source_id VARCHAR(255) NOT NULL,
    format VARCHAR(10) NOT NULL,
    date_from TIMESTAMP WITH TIME ZONE NOT NULL,
    date_to TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    file_url TEXT,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_compliance_exports_source_id ON compliance_exports(source_id);
CREATE INDEX IF NOT EXISTS idx_compliance_exports_status ON compliance_exports(status);

-- Create status_page_configs table if not exists
CREATE TABLE IF NOT EXISTS status_page_configs (
    id SERIAL PRIMARY KEY,
    source_id VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    logo_url TEXT,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_status_page_configs_slug ON status_page_configs(slug);

-- Create status_page_embeds table if not exists
CREATE TABLE IF NOT EXISTS status_page_embeds (
    id SERIAL PRIMARY KEY,
    source_id VARCHAR(255) NOT NULL,
    embed_token VARCHAR(255) NOT NULL UNIQUE,
    theme VARCHAR(20) DEFAULT 'auto',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_status_page_embeds_token ON status_page_embeds(embed_token);

-- Create config_audit_trails table if not exists
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

-- Create config_webhooks table if not exists
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
DROP TABLE IF EXISTS status_page_embeds;
DROP TABLE IF EXISTS status_page_configs;
DROP TABLE IF EXISTS compliance_exports;
DROP TABLE IF EXISTS activity_feeds;
DROP TABLE IF EXISTS sessions;
