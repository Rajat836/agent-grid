-- Materialized view that populates api_metrics from api_executions
CREATE MATERIALIZED VIEW IF NOT EXISTS api_metrics_mv TO api_metrics AS
SELECT
    api_deployments_id AS apis_deployment_id,
    toStartOfFiveMinutes(fromUnixTimestamp(request_timestamp)) AS window_start,
    5 AS window_minutes,
    count() AS total_calls,
    countIf(http_status >= 200 AND http_status < 400) / count() AS success_rate,
    countIf(http_status >= 400) / count() AS error_rate,
    quantile(0.5)(request_duration_ms) AS p50_latency_ms,
    quantile(0.95)(request_duration_ms) AS p95_latency_ms,
    quantile(0.99)(request_duration_ms) AS p99_latency_ms,
    toUnixTimestamp(now()) AS created_at,
    toUnixTimestamp(now()) AS updated_at
FROM api_executions
GROUP BY
    apis_deployment_id,
    window_start,
    window_minutes;
