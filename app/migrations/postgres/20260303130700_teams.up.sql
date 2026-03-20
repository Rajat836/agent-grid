CREATE TABLE
    IF NOT EXISTS teams (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(128) NOT NULL,
        slack_channel VARCHAR(128),
        pagerduty_key VARCHAR(128),
        oncall_email VARCHAR(255),
        created_at TIMESTAMP DEFAULT NOW (),
        updated_at TIMESTAMP DEFAULT NOW (),
        deleted_at TIMESTAMP DEFAULT NULL
    );

CREATE INDEX idx_teams_name ON teams (name);