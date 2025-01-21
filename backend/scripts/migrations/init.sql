-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create ENUM types
CREATE TYPE account_status AS ENUM ('pending', 'active', 'inactive');
CREATE TYPE workspace_status AS ENUM ('active', 'inactive');
CREATE TYPE workspace_role AS ENUM ('admin', 'user', 'viewer');
CREATE TYPE membership_status AS ENUM ('pending', 'active', 'inactive');
CREATE TYPE competitor_status AS ENUM ('active', 'inactive');
CREATE TYPE page_status AS ENUM ('active', 'inactive');
CREATE TYPE history_status AS ENUM ('active', 'inactive');
CREATE TYPE workflow_type AS ENUM ('screenshot', 'reporting');
-- Create workspaces table
CREATE TABLE workspaces (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  billing_email VARCHAR(255) NOT NULL,
  workspace_status workspace_status NOT NULL DEFAULT 'active',
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
  workspace_role workspace_role NOT NULL DEFAULT 'user',
  membership_status membership_status NOT NULL DEFAULT 'pending',
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
  diff_profile TEXT [] DEFAULT ARRAY ['branding', 'customers', 'integration', 'product', 'pricing', 'partnerships', 'messaging'],
  last_checked_at TIMESTAMP WITH TIME ZONE,
  status page_status NOT NULL DEFAULT 'active',
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- Create page_history table
CREATE TABLE page_history (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
  diff_content JSONB,
  prev TEXT,
  current TEXT,
  status history_status NOT NULL DEFAULT 'active',
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE workflow_schedules (
  id UUID PRIMARY KEY,
  workflow_type workflow_type NOT NULL,
  about TEXT,
  spec TEXT NOT NULL,
  last_run TIMESTAMP WITH TIME ZONE,
  next_run TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE TABLE job_records (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  job_id UUID NOT NULL,
  workflow_type TEXT NOT NULL CHECK (workflow_type IN ('screenshot', 'report')),
  start_time TIMESTAMP WITH TIME ZONE,
  end_time TIMESTAMP WITH TIME ZONE,
  cancel_time TIMESTAMP WITH TIME ZONE,
  preemptions INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE
);
-- Create indexes for better query performance
-- Indexes for workspaces
CREATE INDEX idx_workspaces_status ON workspaces(workspace_status);
CREATE INDEX idx_workspaces_slug ON workspaces(slug);
-- Indexes for users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_clerk_id ON users(clerk_id);
CREATE INDEX idx_users_status ON users(status);
-- Indexes for workspace_users
CREATE INDEX idx_workspace_users_workspace_id ON workspace_users(workspace_id);
CREATE INDEX idx_workspace_users_user_id ON workspace_users(user_id);
CREATE INDEX idx_workspace_users_membership_status ON workspace_users(membership_status);
CREATE INDEX idx_workspace_users_workspace_role ON workspace_users(workspace_role);
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
-- Index for listing schedules
CREATE INDEX idx_workflow_schedules_deleted_at ON workflow_schedules(deleted_at NULLS FIRST);
CREATE INDEX idx_workflow_schedules_workflow_type ON workflow_schedules(workflow_type)
WHERE deleted_at IS NULL;
CREATE INDEX idx_job_records_job_id ON job_records(job_id)
WHERE deleted_at IS NULL;
CREATE INDEX idx_job_records_workflow_type ON job_records(workflow_type)
WHERE deleted_at IS NULL;
CREATE INDEX idx_job_records_start_time ON job_records(start_time)
WHERE deleted_at IS NULL;
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