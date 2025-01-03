-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create ENUM types
CREATE TYPE account_status AS ENUM ('pending', 'active', 'inactive');
CREATE TYPE workspace_status AS ENUM ('active', 'inactive');
CREATE TYPE user_workspace_role AS ENUM ('admin', 'user', 'viewer');
CREATE TYPE user_workspace_status AS ENUM ('pending', 'active', 'inactive');
CREATE TYPE competitor_status AS ENUM ('active', 'inactive');
CREATE TYPE page_status AS ENUM ('active', 'inactive');
-- Create workspaces table
CREATE TABLE workspaces (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    billing_email VARCHAR(255) NOT NULL,
    status workspace_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    clerk_id VARCHAR(255) UNIQUE,
    email VARCHAR(255),
    name VARCHAR(255),
    status account_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (email)
);
-- Create workspace_users table (junction table)
CREATE TABLE workspace_users (
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role user_workspace_role NOT NULL DEFAULT 'user',
    status user_workspace_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, user_id)
);
-- Create competitors table
CREATE TABLE competitors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    status competitor_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- Create pages table
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    competitor_id UUID NOT NULL REFERENCES competitors(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    capture_profile JSONB,
    diff_profile JSONB,
    last_checked_at TIMESTAMP WITH TIME ZONE,
    status page_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- Create page_history table
CREATE TABLE page_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    week_number_1 INTEGER NOT NULL,
    week_number_2 INTEGER NOT NULL,
    year_number_1 INTEGER NOT NULL,
    year_number_2 INTEGER NOT NULL,
    bucket_id_1 VARCHAR(255) NOT NULL,
    bucket_id_2 VARCHAR(255) NOT NULL,
    diff_content JSONB,
    screenshot_url_1 TEXT,
    screenshot_url_2 TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- Create indexes for better query performance
-- Indexes for workspaces
CREATE INDEX idx_workspaces_status ON workspaces(status);
CREATE INDEX idx_workspaces_slug ON workspaces(slug);
-- Indexes for users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_clerk_id ON users(clerk_id);
CREATE INDEX idx_users_status ON users(status);
-- Indexes for workspace_users
CREATE INDEX idx_workspace_users_workspace_id ON workspace_users(workspace_id);
CREATE INDEX idx_workspace_users_user_id ON workspace_users(user_id);
CREATE INDEX idx_workspace_users_status ON workspace_users(status);
CREATE INDEX idx_workspace_users_role ON workspace_users(role);
-- Indexes for competitors
CREATE INDEX idx_competitors_workspace_id ON competitors(workspace_id);
CREATE INDEX idx_competitors_status ON competitors(status);
-- Indexes for pages
CREATE INDEX idx_pages_competitor_id ON pages(competitor_id);
CREATE INDEX idx_pages_status ON pages(status);
CREATE INDEX idx_pages_last_checked_at ON pages(last_checked_at);
CREATE INDEX idx_pages_url ON pages(url);
-- Indexes for page_history
CREATE INDEX idx_page_history_page_id ON page_history(page_id);
CREATE INDEX idx_page_history_week_year ON page_history(
    week_number_1,
    year_number_1,
    week_number_2,
    year_number_2
);
-- Functions for updating timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';
-- Create triggers for updating timestamps
CREATE TRIGGER update_workspaces_updated_at BEFORE
UPDATE ON workspaces FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_workspace_users_updated_at BEFORE
UPDATE ON workspace_users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_competitors_updated_at BEFORE
UPDATE ON competitors FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_pages_updated_at BEFORE
UPDATE ON pages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();