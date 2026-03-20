-- Feature-team role mappings derived from Postgres via CDC (Debezium -> Kafka -> ClickHouse)
-- Source of truth remains in Postgres (service_team_roles + features)
-- ReplacingMergeTree allows CDC updates to replace stale rows
CREATE TABLE IF NOT EXISTS feature_team_roles
(
    feature_id Int64,
    team_id Int64,
    role LowCardinality(String),
    assigned_at Int64,

    -- Indexes
    INDEX idx_feature_id feature_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_team_id team_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_role role TYPE set(20) GRANULARITY 1
)
ENGINE = ReplacingMergeTree(assigned_at)
PRIMARY KEY (feature_id, team_id, role)
ORDER BY (feature_id, team_id, role)
SETTINGS index_granularity = 8192;
