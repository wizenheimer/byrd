-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum for user roles
CREATE TYPE user_role AS ENUM ('admin', 'member', 'viewer');

-- Create enum for status
CREATE TYPE record_status AS ENUM ('active', 'deleted');

-- Users table (adapted for Clerk with invitation support)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    clerk_id VARCHAR(255) UNIQUE,           -- Nullable for pending invitations
    email VARCHAR(255) NOT NULL UNIQUE,     -- Primary email from Clerk or invitation
    name VARCHAR(255),                      -- Name from Clerk
    status record_status DEFAULT 'active',  -- Keep status for users to handle account state
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Workspaces table
CREATE TABLE workspaces (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,      -- For unique workspace URLs
    status record_status DEFAULT 'active',  -- Keep status for workspaces to handle deletion
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Workspace Members (Many-to-Many relationship)
CREATE TABLE workspace_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id),
    user_id UUID REFERENCES users(id),
    role user_role NOT NULL DEFAULT 'member',
    status record_status DEFAULT 'active',  -- Keep status for membership management
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(workspace_id, user_id, status)   -- Allow same user to be readded after deletion
);

-- Competitors table
CREATE TABLE competitors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status record_status DEFAULT 'active',  -- Keep status for competitor management
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Pages table (formerly urls)
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    competitor_id UUID REFERENCES competitors(id),
    url TEXT NOT NULL,
    screenshot_options JSONB DEFAULT '{}',
    diff_profile JSONB DEFAULT '{}',
    last_checked_at TIMESTAMP WITH TIME ZONE,
    status record_status DEFAULT 'active',  -- Keep status for page monitoring management
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Diffs table
CREATE TABLE diffs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID REFERENCES pages(id),      -- Updated reference to pages table
    week_number_1 INTEGER NOT NULL,
    week_number_2 INTEGER NOT NULL,
    year_number_1 INTEGER NOT NULL,
    year_number_2 INTEGER NOT NULL,
    bucket_id_1 VARCHAR(100) NOT NULL,
    bucket_id_2 VARCHAR(100) NOT NULL,
    diff_content JSONB NOT NULL,
    screenshot_url_1 TEXT,                  
    screenshot_url_2 TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_week_numbers CHECK (
        week_number_1 BETWEEN 1 AND 53 AND
        week_number_2 BETWEEN 1 AND 53
    ),
    CONSTRAINT valid_year_numbers CHECK (
        year_number_1 > 2000 AND
        year_number_2 > 2000
    )
);