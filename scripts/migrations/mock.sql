-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Enum for user roles in workspaces
CREATE TYPE workspace_role AS ENUM ('owner', 'admin', 'member');
-- Enum for billing plans
CREATE TYPE billing_plan AS ENUM ('trial', 'starter', 'scaleup', 'enterprise');
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Workspaces table
CREATE TABLE workspaces (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    billing_plan billing_plan NOT NULL DEFAULT 'trial',
    max_competitors INT,
    max_urls_per_competitor INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Junction table for users and workspaces with roles
CREATE TABLE workspace_users (
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role workspace_role NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (workspace_id, user_id)
);
-- Competitors table
CREATE TABLE competitors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- URLs table
CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Junction table for competitors and URLs
CREATE TABLE competitor_urls (
    competitor_id UUID REFERENCES competitors(id) ON DELETE CASCADE,
    url_id UUID REFERENCES urls(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    -- Single status field to control monitoring
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (competitor_id, url_id)
);
-- Diffs table to store changes
CREATE TABLE diffs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url_id UUID REFERENCES urls(id) ON DELETE CASCADE,
    detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    changes JSONB NOT NULL,
    diff_type VARCHAR(50),
    -- e.g., 'content', 'visual'
    change_percentage DECIMAL(5, 2),
    -- Stores percentage of change (e.g., 25.50 for 25.5%)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Add indexes for common queries
CREATE INDEX idx_workspace_users_user_id ON workspace_users(user_id);
CREATE INDEX idx_competitors_workspace_id ON competitors(workspace_id);
CREATE INDEX idx_competitor_urls_url_id ON competitor_urls(url_id);
CREATE INDEX idx_diffs_url_id ON diffs(url_id);
CREATE INDEX idx_diffs_detected_at ON diffs(detected_at);
-- Add triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';
CREATE TRIGGER update_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_workspaces_updated_at BEFORE
UPDATE ON workspaces FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_competitors_updated_at BEFORE
UPDATE ON competitors FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_urls_updated_at BEFORE
UPDATE ON urls FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
