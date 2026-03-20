CREATE TABLE IF NOT EXISTS feature_instances
(
    -- UUID generated at app layer
    id UUID,

    -- Cross-DB FK to Postgres features table (enforced at app layer)
    feature_id Int64,

    -- Instance tracking
    status LowCardinality(String),         -- e.g. 'running', 'completed', 'failed'

    -- Timestamps (Unix seconds)
    started_at Int64,
    completed_at Nullable(Int64),
    created_at Int64 DEFAULT toUnixTimestamp(now()),

    -- Indexes
    INDEX idx_feature_id feature_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_status status TYPE set(20) GRANULARITY 1
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(fromUnixTimestamp(started_at))
PRIMARY KEY (started_at, id)
ORDER BY (started_at, id, feature_id)
SETTINGS index_granularity = 8192;
