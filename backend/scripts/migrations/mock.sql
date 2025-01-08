-- Insert workspaces
INSERT INTO workspaces (name, slug, billing_email, status)
VALUES (
        'Acme Corporation',
        'acme-corp',
        'billing@acme.com',
        'active'
    ),
    (
        'TechStart Inc',
        'techstart',
        'finance@techstart.io',
        'active'
    ),
    (
        'Global Systems',
        'global-sys',
        'accounts@globalsys.com',
        'active'
    ),
    (
        'Dev Labs',
        'dev-labs',
        'billing@devlabs.tech',
        'inactive'
    );
-- Insert users
INSERT INTO users (clerk_id, email, name, status)
VALUES (
        'clk_123',
        'john.doe@example.com',
        'John Doe',
        'active'
    ),
    (
        'clk_124',
        'jane.smith@example.com',
        'Jane Smith',
        'active'
    ),
    (
        'clk_125',
        'bob.wilson@example.com',
        'Bob Wilson',
        'active'
    ),
    (
        'clk_126',
        'alice.cooper@example.com',
        'Alice Cooper',
        'pending'
    ),
    (
        'clk_127',
        'charlie.brown@example.com',
        'Charlie Brown',
        'active'
    ),
    (
        NULL,
        'pending.user@example.com',
        'Pending User',
        'pending'
    ),
    (
        NULL,
        'invited.user@example.com',
        'Invited User',
        'pending'
    );
-- Insert workspace_users (using subqueries to get IDs)
INSERT INTO workspace_users (workspace_id, user_id, role, status)
SELECT w.id,
    u.id,
    role_status.role::user_workspace_role,
    role_status.status::user_workspace_status
FROM (
        SELECT 'Acme Corporation' as workspace_name,
            'john.doe@example.com' as email,
            'admin' as role,
            'active' as status
        UNION ALL
        SELECT 'Acme Corporation',
            'jane.smith@example.com',
            'user',
            'active'
        UNION ALL
        SELECT 'TechStart Inc',
            'bob.wilson@example.com',
            'admin',
            'active'
        UNION ALL
        SELECT 'TechStart Inc',
            'alice.cooper@example.com',
            'viewer',
            'pending'
        UNION ALL
        SELECT 'Global Systems',
            'charlie.brown@example.com',
            'admin',
            'active'
        UNION ALL
        SELECT 'Global Systems',
            'pending.user@example.com',
            'user',
            'pending'
    ) as role_status
    JOIN workspaces w ON w.name = role_status.workspace_name
    JOIN users u ON u.email = role_status.email;
-- Insert competitors
INSERT INTO competitors (workspace_id, name, status)
SELECT w.id,
    comp.name,
    comp.status::competitor_status
FROM (
        SELECT 'Acme Corporation' as workspace_name,
            'Competitor A' as name,
            'active' as status
        UNION ALL
        SELECT 'Acme Corporation',
            'Competitor B',
            'active'
        UNION ALL
        SELECT 'TechStart Inc',
            'Competitor C',
            'active'
        UNION ALL
        SELECT 'TechStart Inc',
            'Competitor D',
            'inactive'
        UNION ALL
        SELECT 'Global Systems',
            'Competitor E',
            'active'
    ) as comp
    JOIN workspaces w ON w.name = comp.workspace_name;
-- Insert pages
INSERT INTO pages (
        competitor_id,
        url,
        capture_profile,
        diff_profile,
        last_checked_at,
        status
    )
SELECT c.id,
    page_data.url,
    page_data.capture_profile::jsonb,
    page_data.diff_profile::jsonb,
    CASE
        WHEN page_data.status = 'active' THEN NOW() - (random() * interval '7 days')
        ELSE NULL
    END as last_checked_at,
    page_data.status::page_status
FROM (
        SELECT 'Competitor A' as competitor_name,
            'https://competitor-a.com/products' as url,
            '{"viewport": {"width": 1920, "height": 1080}, "waitUntil": "networkidle0"}' as capture_profile,
            '{"threshold": 0.1, "ignoreSelectors": [".ads", ".dynamic-content"]}' as diff_profile,
            'active' as status
        UNION ALL
        SELECT 'Competitor A',
            'https://competitor-a.com/pricing',
            '{"viewport": {"width": 1920, "height": 1080}}',
            '{"threshold": 0.2}',
            'active'
        UNION ALL
        SELECT 'Competitor B',
            'https://competitor-b.com/features',
            '{"viewport": {"width": 1366, "height": 768}}',
            '{"threshold": 0.15}',
            'active'
        UNION ALL
        SELECT 'Competitor C',
            'https://competitor-c.com/about',
            '{"viewport": {"width": 1440, "height": 900}}',
            '{"threshold": 0.1}',
            'inactive'
    ) as page_data
    JOIN competitors c ON c.name = page_data.competitor_name;
-- Insert page_history
INSERT INTO page_history (
        page_id,
        week_number_1,
        week_number_2,
        year_number_1,
        year_number_2,
        bucket_id_1,
        bucket_id_2,
        diff_content,
        screenshot_url_1,
        screenshot_url_2
    )
SELECT p.id,
    EXTRACT(
        WEEK
        FROM CURRENT_DATE - interval '1 week'
    )::integer,
    EXTRACT(
        WEEK
        FROM CURRENT_DATE
    )::integer,
    EXTRACT(
        YEAR
        FROM CURRENT_DATE - interval '1 week'
    )::integer,
    EXTRACT(
        YEAR
        FROM CURRENT_DATE
    )::integer,
    'bucket_' || (random() * 1000)::integer,
    'bucket_' || (random() * 1000)::integer,
    '{"changes": ["header.logo", "pricing.table"], "similarity": 0.85}'::jsonb,
    'https://storage.example.com/screenshots/' || p.id || '/week_' || EXTRACT(
        WEEK
        FROM CURRENT_DATE - interval '1 week'
    )::integer || '.png',
    'https://storage.example.com/screenshots/' || p.id || '/week_' || EXTRACT(
        WEEK
        FROM CURRENT_DATE
    )::integer || '.png'
FROM pages p
WHERE p.status = 'active'
    AND random() < 0.7;
-- Only create history for some pages
-- Insert additional history entries for some pages (older entries)
INSERT INTO page_history (
        page_id,
        week_number_1,
        week_number_2,
        year_number_1,
        year_number_2,
        bucket_id_1,
        bucket_id_2,
        diff_content,
        screenshot_url_1,
        screenshot_url_2
    )
SELECT p.id,
    EXTRACT(
        WEEK
        FROM CURRENT_DATE - interval '2 weeks'
    )::integer,
    EXTRACT(
        WEEK
        FROM CURRENT_DATE - interval '1 week'
    )::integer,
    EXTRACT(
        YEAR
        FROM CURRENT_DATE - interval '2 weeks'
    )::integer,
    EXTRACT(
        YEAR
        FROM CURRENT_DATE - interval '1 week'
    )::integer,
    'bucket_' || (random() * 1000)::integer,
    'bucket_' || (random() * 1000)::integer,
    '{"changes": ["footer.links", "content.main"], "similarity": 0.92}'::jsonb,
    'https://storage.example.com/screenshots/' || p.id || '/week_' || EXTRACT(
        WEEK
        FROM CURRENT_DATE - interval '2 weeks'
    )::integer || '.png',
    'https://storage.example.com/screenshots/' || p.id || '/week_' || EXTRACT(
        WEEK
        FROM CURRENT_DATE - interval '1 week'
    )::integer || '.png'
FROM pages p
WHERE p.status = 'active'
    AND random() < 0.5;
-- Only create older history for some pages