CREATE TABLE IF NOT EXISTS api_db_interactions
(
    -- UUID generated at app layer
    id UUID,

    -- FK to ClickHouse api_executions table
    api_execution_id UUID,

    -- Interaction details
    database_type LowCardinality(String),      -- e.g. 'postgres', 'clickhouse', 'redis'
    table_name String,
    operation_type LowCardinality(String),      -- e.g. 'SELECT', 'INSERT', 'UPDATE', 'DELETE'
    sql_query String,
    query_parameters Nullable(String),

    -- Performance / result
    rows_affected Int32 DEFAULT 0,
    duration_ms Int32 DEFAULT 0,
    error_message Nullable(String),

    -- Metadata
    created_at Int64 DEFAULT toUnixTimestamp(now()),

    -- Indexes
    INDEX idx_api_execution_id api_execution_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_database_type database_type TYPE set(10) GRANULARITY 1,
    INDEX idx_table_name table_name TYPE bloom_filter GRANULARITY 1,
    INDEX idx_operation_type operation_type TYPE set(10) GRANULARITY 1,
    INDEX idx_duration_ms duration_ms TYPE minmax GRANULARITY 1
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(fromUnixTimestamp(created_at))
PRIMARY KEY (created_at, id)
ORDER BY (created_at, id, api_execution_id)
SETTINGS index_granularity = 8192;
