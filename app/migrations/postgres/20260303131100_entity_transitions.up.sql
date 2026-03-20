CREATE TABLE IF NOT EXISTS entity_transitions (
    id BIGSERIAL PRIMARY KEY,
    from_entity_id BIGINT NOT NULL,
    to_entity_id BIGINT NOT NULL,
    condition_description TEXT,          -- human-readable, not executable logic (e.g. 'amount < 200')
    condition_expression TEXT,           -- executable logic (e.g. 'amount < 200')
    transition_type VARCHAR(64),         -- e.g. 'conditional', 'automatic', 'manual'
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_entity_transitions_from_entity_id ON entity_transitions(from_entity_id);
CREATE INDEX idx_entity_transitions_to_entity_id ON entity_transitions(to_entity_id);
CREATE INDEX idx_entity_transitions_transition_type ON entity_transitions(transition_type);
