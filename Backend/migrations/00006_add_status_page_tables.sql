-- +goose Up
-- Create status_page_configs table
CREATE TABLE IF NOT EXISTS status_page_configs (
    id SERIAL PRIMARY KEY,
    source_id VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    logo_url TEXT,
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_status_page_configs_slug ON status_page_configs(slug);
CREATE INDEX IF NOT EXISTS idx_status_page_configs_source_id ON status_page_configs(source_id);

-- Create status_page_embeds table
CREATE TABLE IF NOT EXISTS status_page_embeds (
    id SERIAL PRIMARY KEY,
    source_id VARCHAR(255) NOT NULL,
    embed_token VARCHAR(255) NOT NULL UNIQUE,
    theme VARCHAR(20) DEFAULT 'auto',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_status_page_embeds_token ON status_page_embeds(embed_token);
CREATE INDEX IF NOT EXISTS idx_status_page_embeds_source_id ON status_page_embeds(source_id);

-- +goose Down
DROP TABLE IF EXISTS status_page_embeds;
DROP TABLE IF EXISTS status_page_configs;
