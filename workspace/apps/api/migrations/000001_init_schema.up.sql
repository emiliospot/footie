-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For text search and similarity

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    avatar TEXT,
    organization VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT true,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Create teams table
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    short_name VARCHAR(100) NOT NULL,
    code VARCHAR(10) UNIQUE NOT NULL,
    country VARCHAR(100) NOT NULL,
    city VARCHAR(100),
    stadium VARCHAR(255),
    stadium_capacity INTEGER,
    founded INTEGER,
    logo TEXT,
    colors VARCHAR(100),
    website VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Create players table
CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    date_of_birth DATE,
    nationality VARCHAR(100),
    position VARCHAR(50) NOT NULL,
    shirt_number INTEGER,
    height INTEGER, -- in cm
    weight INTEGER, -- in kg
    preferred_foot VARCHAR(10),
    photo TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Create matches table
CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    home_team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    away_team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    match_date TIMESTAMPTZ NOT NULL,
    competition VARCHAR(100) NOT NULL,
    season VARCHAR(20) NOT NULL,
    round VARCHAR(50),
    stadium VARCHAR(255),
    attendance INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled', -- scheduled, live, finished, postponed, canceled
    referee VARCHAR(255),
    home_team_score INTEGER NOT NULL DEFAULT 0,
    away_team_score INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT different_teams CHECK (home_team_id != away_team_id)
);

-- Create match_events table (for goals, cards, substitutions, shots, passes, etc.)
CREATE TABLE match_events (
    id SERIAL PRIMARY KEY,
    match_id INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    team_id INTEGER REFERENCES teams(id) ON DELETE SET NULL,
    player_id INTEGER REFERENCES players(id) ON DELETE SET NULL,
    secondary_player_id INTEGER REFERENCES players(id) ON DELETE SET NULL, -- for assists, substitutions
    event_type VARCHAR(50) NOT NULL, -- goal, yellow_card, red_card, substitution, shot, pass, tackle, etc.
    minute INTEGER NOT NULL,
    extra_minute INTEGER DEFAULT 0,
    position_x NUMERIC(5,2), -- field position coordinates for analytics
    position_y NUMERIC(5,2),
    description TEXT,
    metadata JSONB, -- for storing additional event data (xG, pass completion, etc.)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Create player_statistics table (aggregated stats per season/competition)
CREATE TABLE player_statistics (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    season VARCHAR(20) NOT NULL,
    competition VARCHAR(100) NOT NULL,

    -- Appearance stats
    matches_played INTEGER NOT NULL DEFAULT 0,
    matches_started INTEGER NOT NULL DEFAULT 0,
    minutes_played INTEGER NOT NULL DEFAULT 0,
    sub_on INTEGER NOT NULL DEFAULT 0,
    sub_off INTEGER NOT NULL DEFAULT 0,

    -- Attacking stats
    goals INTEGER NOT NULL DEFAULT 0,
    assists INTEGER NOT NULL DEFAULT 0,
    shots_total INTEGER NOT NULL DEFAULT 0,
    shots_on_target INTEGER NOT NULL DEFAULT 0,
    shot_accuracy NUMERIC(5,2) DEFAULT 0,
    goal_conversion NUMERIC(5,2) DEFAULT 0,

    -- Passing stats
    passes_total INTEGER NOT NULL DEFAULT 0,
    passes_completed INTEGER NOT NULL DEFAULT 0,
    pass_accuracy NUMERIC(5,2) DEFAULT 0,
    key_passes INTEGER NOT NULL DEFAULT 0,
    crosses INTEGER NOT NULL DEFAULT 0,

    -- Defensive stats
    tackles INTEGER NOT NULL DEFAULT 0,
    tackles_won INTEGER NOT NULL DEFAULT 0,
    interceptions INTEGER NOT NULL DEFAULT 0,
    clearances INTEGER NOT NULL DEFAULT 0,
    blocked_shots INTEGER NOT NULL DEFAULT 0,

    -- Duels
    duels INTEGER NOT NULL DEFAULT 0,
    duels_won INTEGER NOT NULL DEFAULT 0,
    aerial_duels INTEGER NOT NULL DEFAULT 0,
    aerial_duels_won INTEGER NOT NULL DEFAULT 0,

    -- Discipline
    yellow_cards INTEGER NOT NULL DEFAULT 0,
    red_cards INTEGER NOT NULL DEFAULT 0,
    fouls INTEGER NOT NULL DEFAULT 0,
    fouls_drawn INTEGER NOT NULL DEFAULT 0,

    -- Goalkeeper stats (nullable for outfield players)
    clean_sheets INTEGER,
    goals_conceded INTEGER,
    saves_total INTEGER,
    save_percentage NUMERIC(5,2),
    penalties_saved INTEGER,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    UNIQUE(player_id, season, competition)
);

-- Create team_statistics table (aggregated stats per season/competition)
CREATE TABLE team_statistics (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    season VARCHAR(20) NOT NULL,
    competition VARCHAR(100) NOT NULL,

    -- Match results
    matches_played INTEGER NOT NULL DEFAULT 0,
    wins INTEGER NOT NULL DEFAULT 0,
    draws INTEGER NOT NULL DEFAULT 0,
    losses INTEGER NOT NULL DEFAULT 0,
    points INTEGER NOT NULL DEFAULT 0,
    position INTEGER,

    -- Goals
    goals_scored INTEGER NOT NULL DEFAULT 0,
    goals_conceded INTEGER NOT NULL DEFAULT 0,
    goal_difference INTEGER NOT NULL DEFAULT 0,
    clean_sheets INTEGER NOT NULL DEFAULT 0,
    goals_per_match NUMERIC(5,2) DEFAULT 0,

    -- Home/Away splits
    home_wins INTEGER NOT NULL DEFAULT 0,
    home_draws INTEGER NOT NULL DEFAULT 0,
    home_losses INTEGER NOT NULL DEFAULT 0,
    away_wins INTEGER NOT NULL DEFAULT 0,
    away_draws INTEGER NOT NULL DEFAULT 0,
    away_losses INTEGER NOT NULL DEFAULT 0,

    -- Team stats
    possession NUMERIC(5,2) DEFAULT 0,
    pass_accuracy NUMERIC(5,2) DEFAULT 0,
    shots_per_match NUMERIC(5,2) DEFAULT 0,
    shots_on_target_percentage NUMERIC(5,2) DEFAULT 0,

    -- Discipline
    yellow_cards INTEGER NOT NULL DEFAULT 0,
    red_cards INTEGER NOT NULL DEFAULT 0,

    -- Form (last 5 matches: W/D/L)
    current_form VARCHAR(10),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    UNIQUE(team_id, season, competition)
);

-- Create indexes for better query performance (critical for analytics)

-- Users indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;

-- Teams indexes
CREATE INDEX idx_teams_name ON teams(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_teams_code ON teams(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_teams_country ON teams(country) WHERE deleted_at IS NULL;
CREATE INDEX idx_teams_name_trgm ON teams USING gin(name gin_trgm_ops) WHERE deleted_at IS NULL; -- For fuzzy search

-- Players indexes
CREATE INDEX idx_players_team_id ON players(team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_players_full_name ON players(full_name) WHERE deleted_at IS NULL;
CREATE INDEX idx_players_position ON players(position) WHERE deleted_at IS NULL;
CREATE INDEX idx_players_nationality ON players(nationality) WHERE deleted_at IS NULL;
CREATE INDEX idx_players_name_trgm ON players USING gin(full_name gin_trgm_ops) WHERE deleted_at IS NULL; -- For fuzzy search

-- Matches indexes
CREATE INDEX idx_matches_home_team ON matches(home_team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_away_team ON matches(away_team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_date ON matches(match_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_competition ON matches(competition) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_season ON matches(season) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_status ON matches(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_comp_season ON matches(competition, season) WHERE deleted_at IS NULL;

-- Match events indexes (critical for analytics queries)
CREATE INDEX idx_match_events_match_id ON match_events(match_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_match_events_player_id ON match_events(player_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_match_events_team_id ON match_events(team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_match_events_type ON match_events(event_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_match_events_minute ON match_events(minute) WHERE deleted_at IS NULL;
CREATE INDEX idx_match_events_metadata ON match_events USING gin(metadata) WHERE deleted_at IS NULL; -- For JSONB queries

-- Player statistics indexes
CREATE INDEX idx_player_stats_player_id ON player_statistics(player_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_player_stats_season ON player_statistics(season) WHERE deleted_at IS NULL;
CREATE INDEX idx_player_stats_competition ON player_statistics(competition) WHERE deleted_at IS NULL;
CREATE INDEX idx_player_stats_goals ON player_statistics(goals DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_player_stats_assists ON player_statistics(assists DESC) WHERE deleted_at IS NULL;

-- Team statistics indexes
CREATE INDEX idx_team_stats_team_id ON team_statistics(team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_team_stats_season ON team_statistics(season) WHERE deleted_at IS NULL;
CREATE INDEX idx_team_stats_competition ON team_statistics(competition) WHERE deleted_at IS NULL;
CREATE INDEX idx_team_stats_points ON team_statistics(points DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_team_stats_position ON team_statistics(position) WHERE deleted_at IS NULL;

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_teams_updated_at BEFORE UPDATE ON teams
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_players_updated_at BEFORE UPDATE ON players
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_matches_updated_at BEFORE UPDATE ON matches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_match_events_updated_at BEFORE UPDATE ON match_events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_player_statistics_updated_at BEFORE UPDATE ON player_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_team_statistics_updated_at BEFORE UPDATE ON team_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
