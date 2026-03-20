CREATE TABLE IF NOT EXISTS service_deployments
(
    id Int64,

    -- Cross-DB FK to Postgres services table (enforced at app layer)
    service_id Int64,

    -- Deployment details
    environment LowCardinality(String),    -- e.g. 'production', 'staging', 'development'
    commit_hash String,                    -- Git commit hash for deployment
    version String,

    -- Timestamps (Unix seconds)
    deployed_at Int64,
    created_at Int64 DEFAULT toUnixTimestamp(now()),
    updated_at Int64 DEFAULT toUnixTimestamp(now()),

    -- Indexes
    INDEX idx_service_id service_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_environment environment TYPE set(20) GRANULARITY 1,
    INDEX idx_version version TYPE bloom_filter GRANULARITY 1,
    INDEX idx_commit_hash commit_hash TYPE bloom_filter GRANULARITY 1
)
ENGINE = ReplacingMergeTree(updated_at)
PARTITION BY toYYYYMM(fromUnixTimestamp(deployed_at))
PRIMARY KEY (deployed_at, id)
ORDER BY (deployed_at, id, service_id)
SETTINGS index_granularity = 8192;