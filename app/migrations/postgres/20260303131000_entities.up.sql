CREATE TABLE IF NOT EXISTS entities (
    id BIGSERIAL PRIMARY KEY,
    feature_id BIGINT,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    display_order INTEGER DEFAULT 0,
    is_start BOOLEAN DEFAULT FALSE,
    is_terminal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_entities_feature_id ON entities(feature_id);
CREATE INDEX idx_entities_code ON entities(code);
CREATE INDEX idx_entities_name ON entities(name);
CREATE INDEX idx_entities_is_start ON entities(is_start);
CREATE INDEX idx_entities_is_terminal ON entities(is_terminal);
