-- Materialized view that populates feature_metrics from feature_instances
CREATE MATERIALIZED VIEW IF NOT EXISTS feature_metrics_mv TO feature_metrics AS
SELECT
    feature_id,
    toStartOfFiveMinutes(fromUnixTimestamp(started_at)) AS window_start,
    5 AS window_minutes,
    count() AS total_count,
    countIf(status = 'completed') / count() AS success_rate,
    countIf(status = 'failed') / count() AS failure_rate,
    quantile(0.5)(if(completed_at IS NOT NULL, completed_at - started_at, 0) * 1000) AS p50_duration_ms,
    quantile(0.95)(if(completed_at IS NOT NULL, completed_at - started_at, 0) * 1000) AS p95_duration_ms,
    quantile(0.99)(if(completed_at IS NOT NULL, completed_at - started_at, 0) * 1000) AS p99_duration_ms,
    toUnixTimestamp(now()) AS created_at
FROM feature_instances
GROUP BY
    feature_id,
    window_start,
    window_minutes;
