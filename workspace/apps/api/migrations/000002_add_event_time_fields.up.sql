-- Add second and period fields to match_events table
-- This allows precise timing (exact second) and period tracking (1st half, 2nd half, extra time, penalties)

ALTER TABLE match_events
ADD COLUMN second INTEGER DEFAULT 0 CHECK (second >= 0 AND second < 60), -- Exact second within the minute (0-59)
ADD COLUMN period VARCHAR(20) DEFAULT 'regular'; -- Period: 'first_half', 'second_half', 'extra_time_first', 'extra_time_second', 'penalties'

-- Create index for period-based queries
CREATE INDEX idx_match_events_period ON match_events(period) WHERE deleted_at IS NULL;

-- Create composite index for time-based queries (minute + second)
CREATE INDEX idx_match_events_time ON match_events(match_id, period, minute, second) WHERE deleted_at IS NULL;

-- Update existing records to have default period based on minute
UPDATE match_events
SET period = CASE
    WHEN minute <= 45 THEN 'first_half'
    WHEN minute <= 90 THEN 'second_half'
    WHEN minute <= 105 THEN 'extra_time_first'
    WHEN minute <= 120 THEN 'extra_time_second'
    ELSE 'penalties'
END
WHERE period = 'regular';
