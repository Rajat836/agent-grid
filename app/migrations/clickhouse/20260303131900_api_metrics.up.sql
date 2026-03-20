-- Target table for the api_metrics materialized view
-- Aggregated from api_executions in 5-minute windows
-- Deduplication via ReplacingMergeTree on updated_at
CREATE TABLE IF NOT EXISTS api_metrics
(
    apis_deployment_id UUID,
    window_start DateTime,
    window_minutes Int32 DEFAULT 5,

    -- Aggregated metrics
    total_calls Int64 DEFAULT 0,
    success_rate Float64 DEFAULT 0,
    error_rate Float64 DEFAULT 0,
    p50_latency_ms Float64 DEFAULT 0,
    p95_latency_ms Float64 DEFAULT 0,
    p99_latency_ms Float64 DEFAULT 0,

    -- Metadata
    created_at Int64 DEFAULT toUnixTimestamp(now()),
    updated_at Int64 DEFAULT toUnixTimestamp(now())
)
ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY toYYYYMM(window_start)
PRIMARY KEY (apis_deployment_id, window_start, window_minutes)
ORDER BY (apis_deployment_id, window_start, window_minutes)
SETTINGS index_granularity = 8192;
