-- Updated Schema
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Create indices
CREATE INDEX idx_urls_id ON urls(id DESC);
CREATE INDEX idx_urls_created_at ON urls(created_at DESC);
-- Insert dummy data
INSERT INTO urls (url)
VALUES ('https://www.google.com'),
    ('https://www.github.com'),
    ('https://www.amazon.com'),
    ('https://www.netflix.com'),
    ('https://docs.docker.com'),
    ('https://kubernetes.io/docs'),
    ('https://go.dev/learn'),
    ('https://www.rust-lang.org/learn'),
    ('https://www.typescriptlang.org/docs/handbook'),
    ('https://www.ruby-lang.org/en/documentation/'),
    ('https://www.php.net/docs.php'),
    ('https://www.python.org/doc'),
    ('https://www.postgresql.org/docs'),
    ('https://www.rust-lang.org/');