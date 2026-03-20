CREATE TABLE
    IF NOT EXISTS services (
        id BIGSERIAL PRIMARY KEY,
        code varchar(128) NOT NULL UNIQUE,
        name VARCHAR(128) NOT NULL,
        description TEXT,
        repository_url VARCHAR(512),
        criticality_level varchar(128), -- e.g. 'critical', 'high', 'medium', 'low'
        created_at TIMESTAMP DEFAULT NOW (),
        updated_at TIMESTAMP DEFAULT NOW (),
        deleted_at TIMESTAMP DEFAULT NULL
    );

CREATE INDEX idx_services_code ON services (code);

CREATE INDEX idx_services_name ON services (name);

CREATE INDEX idx_services_criticality_level ON services (criticality_level);