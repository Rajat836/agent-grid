CREATE TABLE IF NOT EXISTS entity_apis (
    id BIGSERIAL PRIMARY KEY,
    entity_id BIGINT NOT NULL,
    api_id BIGINT NOT NULL
);

CREATE INDEX idx_entity_apis_entity_id ON entity_apis(entity_id);
CREATE INDEX idx_entity_apis_api_id ON entity_apis(api_id);
