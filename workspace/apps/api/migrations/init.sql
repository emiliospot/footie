-- Initial database setup
-- This file is executed when the PostgreSQL container first starts

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For text search

-- Create indexes that might not be automatically created by GORM
-- Additional indexes can be added here for performance optimization

-- Example: Create GIN index for JSONB columns (if needed)
-- CREATE INDEX IF NOT EXISTS idx_match_events_metadata ON match_events USING gin(metadata);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Note: Triggers will be added after tables are created by GORM migrations

