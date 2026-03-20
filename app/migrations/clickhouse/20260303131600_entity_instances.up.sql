CREATE TABLE IF NOT EXISTS entity_instances
(
    -- UUID generated at app layer
    id UUID,

    -- FK to ClickHouse feature_instances table
    feature_instance_id UUID,

    -- Cross-DB FK to Postgres entities table (enforced at app layer)
    entity_id Int64,

    -- Instance tracking
    status LowCardinality(String),         -- e.g. 'pending', 'active', 'completed', 'failed'

    -- Timestamps (Unix seconds)
    started_at Int64,
    completed_at Nullable(Int64),
    created_at Int64 DEFAULT toUnixTimestamp(now()),

    -- Indexes
    INDEX idx_feature_instance_id feature_instance_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_entity_id entity_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_status status TYPE set(20) GRANULARITY 1
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(fromUnixTimestamp(started_at))
PRIMARY KEY (started_at, id)
ORDER BY (started_at, id, feature_instance_id)
SETTINGS index_granularity = 8192;
