CREATE TABLE IF NOT EXISTS apis_deployments
(
    -- UUID generated at app layer
    id UUID,

    -- Cross-DB FK to Postgres apis table (enforced at app layer)
    api_id Int64,

    -- FK to ClickHouse service_deployments table
    service_deployment_id Int64,

    -- Indexes
    INDEX idx_api_id api_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_service_deployment_id service_deployment_id TYPE bloom_filter GRANULARITY 1
)
ENGINE = MergeTree()
PRIMARY KEY (id)
ORDER BY (id, api_id, service_deployment_id)
SETTINGS index_granularity = 8192;
