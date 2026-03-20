CREATE TABLE IF NOT EXISTS kpis_relationships (
    id BIGSERIAL PRIMARY KEY,
    kpi_id BIGINT NOT NULL,
    relation_type VARCHAR(64) NOT NULL,        -- e.g. 'contributes_to', 'depends_on'
    target_type VARCHAR(64) NOT NULL,          -- polymorphic: 'service', 'team', 'feature', etc.
    target_id BIGINT NOT NULL,
    weight FLOAT DEFAULT 0,
    weight_set_by VARCHAR(64),                 -- user/system that set the weight
    weight_reviewed_by VARCHAR(64),            -- reviewer identifier
    weight_reviewed_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_kpis_relationships_kpi_id ON kpis_relationships(kpi_id);
CREATE INDEX idx_kpis_relationships_target ON kpis_relationships(target_type, target_id);
CREATE INDEX idx_kpis_relationships_relation_type ON kpis_relationships(relation_type);
