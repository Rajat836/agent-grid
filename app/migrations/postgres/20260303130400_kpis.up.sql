CREATE TABLE IF NOT EXISTS kpis (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    metric_type VARCHAR(64) NOT NULL,  -- e.g. 'ratio', 'count', 'percentage', 'duration'
    unit VARCHAR(32),                  -- e.g. 'ms', '%', 'count'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_kpis_code ON kpis(code);
CREATE INDEX idx_kpis_name ON kpis(name);
CREATE INDEX idx_kpis_metric_type ON kpis(metric_type);
