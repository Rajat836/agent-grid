CREATE TABLE
    IF NOT EXISTS kpis (
        id BIGSERIAL PRIMARY KEY,
        code varchar(128) NOT NULL UNIQUE,
        name VARCHAR(128) NOT NULL,
        description TEXT,
        metric_type VARCHAR(128) NOT NULL, -- e.g. 'ratio', 'count', 'percentage', 'duration'
        unit varchar(128), -- e.g. 'ms', '%', 'count'
        created_at TIMESTAMP DEFAULT NOW (),
        updated_at TIMESTAMP DEFAULT NOW ()
    );

CREATE INDEX idx_kpis_code ON kpis (code);

CREATE INDEX idx_kpis_name ON kpis (name);

CREATE INDEX idx_kpis_metric_type ON kpis (metric_type);