-- Create indexes for better query performance
CREATE INDEX idx_users_clerk_id ON users(clerk_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_workspace_members_user_id ON workspace_members(user_id);
CREATE INDEX idx_workspace_members_workspace_id ON workspace_members(workspace_id);
CREATE INDEX idx_workspace_members_status ON workspace_members(status);
CREATE INDEX idx_competitors_workspace_id ON competitors(workspace_id);
CREATE INDEX idx_competitors_status ON competitors(status);
CREATE INDEX idx_pages_competitor_id ON pages(competitor_id);
CREATE INDEX idx_pages_status ON pages(status);
CREATE INDEX idx_diffs_page_id ON diffs(page_id);
CREATE INDEX idx_pages_last_checked_at ON pages(last_checked_at);
CREATE INDEX idx_diffs_week_year ON diffs(week_number_1, year_number_1, week_number_2, year_number_2);