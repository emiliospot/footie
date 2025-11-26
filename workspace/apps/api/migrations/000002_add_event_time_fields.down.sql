-- Remove second and period fields
DROP INDEX IF EXISTS idx_match_events_period;
DROP INDEX IF EXISTS idx_match_events_time;

ALTER TABLE match_events
DROP COLUMN IF EXISTS second,
DROP COLUMN IF EXISTS period;

