-- Target table for the entity_metrics materialized view
-- Aggregated from entity_instances in 5-minute windows
CREATE TABLE IF NOT EXISTS entity_metrics
(
    entity_id Int64,
    window_start DateTime,
    window_minutes Int32 DEFAULT 5,

    -- Aggregated metrics
    total_count Int64 DEFAULT 0,
    success_rate Float64 DEFAULT 0,
    failure_rate Float64 DEFAULT 0,
    p50_duration_ms Float64 DEFAULT 0,
    p95_duration_ms Float64 DEFAULT 0,
    p99_duration_ms Float64 DEFAULT 0,

    -- Metadata
    created_at Int64 DEFAULT toUnixTimestamp(now())
)
ENGINE = ReplacingMergeTree(created_at)
PARTITION BY toYYYYMM(window_start)
PRIMARY KEY (entity_id, window_start, window_minutes)
ORDER BY (entity_id, window_start, window_minutes)
SETTINGS index_granularity = 8192;
