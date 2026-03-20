CREATE TABLE IF NOT EXISTS api_executions
(
    -- UUID generated at app layer
    id UUID,

    -- FK references (enforced at app layer)
    entity_instance_id UUID,
    api_deployments_id UUID,

    -- Identity / auth context
    oauth_client_id Nullable(String),
    user_id Nullable(String),
    user_uid Nullable(String),
    organization_id Nullable(String),
    token_id Nullable(String),

    -- Request metadata
    request_timestamp Int64,
    request_duration_ms Int32,
    source_ip Nullable(String),
    request_id Nullable(String),
    correlation_id Nullable(String),
    trace_id Nullable(String),
    span_id Nullable(String),
    host Nullable(String),
    latitude Nullable(Float64),
    longitude Nullable(Float64),

    -- Device context
    device_id Nullable(String),
    platform_os LowCardinality(Nullable(String)),
    app_version Nullable(String),
    user_agent Nullable(String),
    risk_session_id Nullable(String),
    device_model Nullable(String),

    -- Auth details
    auth_scheme LowCardinality(Nullable(String)),   -- e.g. 'Bearer', 'Basic', 'API-Key'
    token_type LowCardinality(Nullable(String)),
    token_scopes Nullable(String),

    -- Response details
    http_status Int32,
    app_error_code Nullable(String),
    error_message Nullable(String),
    response_size_bytes Int32 DEFAULT 0,

    -- Metadata
    created_at Int64 DEFAULT toUnixTimestamp(now()),

    -- Indexes
    INDEX idx_entity_instance_id entity_instance_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_api_deployments_id api_deployments_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_user_id user_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_organization_id organization_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_request_id request_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_correlation_id correlation_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_trace_id trace_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_http_status http_status TYPE set(100) GRANULARITY 1,
    INDEX idx_device_id device_id TYPE bloom_filter GRANULARITY 1
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(fromUnixTimestamp(request_timestamp))
PRIMARY KEY (request_timestamp, id)
ORDER BY (request_timestamp, id, entity_instance_id)
SETTINGS index_granularity = 8192;
