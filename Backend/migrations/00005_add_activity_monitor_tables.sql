-- +goose Up
-- Create activity_feeds table
CREATE TABLE IF NOT EXISTS activity_feeds (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    source_id VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(255),
    context JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_activity_feeds_user_id ON activity_feeds(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_feeds_source_id ON activity_feeds(source_id);
CREATE INDEX IF NOT EXISTS idx_activity_feeds_action ON activity_feeds(action);
CREATE INDEX IF NOT EXISTS idx_activity_feeds_created_at ON activity_feeds(created_at DESC);

-- Create compliance_exports table
CREATE TABLE IF NOT EXISTS compliance_exports (
    id SERIAL PRIMARY KEY,
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
CREATE INDEX IF NOT EXISTS idx_compliance_exports_created_at ON compliance_exports(created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS activity_feeds;
DROP TABLE IF EXISTS compliance_exports;
