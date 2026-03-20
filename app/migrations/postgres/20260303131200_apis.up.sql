CREATE TABLE IF NOT EXISTS apis (
    id BIGSERIAL PRIMARY KEY,
    endpoint VARCHAR(512) NOT NULL,
    http_method VARCHAR(16) NOT NULL,    -- e.g. 'GET', 'POST', 'PUT', 'DELETE', 'PATCH'
    is_internal BOOLEAN DEFAULT FALSE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_apis_is_internal ON apis(is_internal);