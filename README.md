# otel-example

This project is a simple example of using OpenTelemetry to collect metrics and sending it to Prometheus and ClickHouse.

```promql
{__name__=~"wave_.*"}
```

```promql
wave_1
```

```promql
sum(wave_1)
```

```promql
avg(wave_1)
```

```promql
wave_1[1h]
```

```SQL
SELECT
MetricName, ResourceAttributes['service.name'] as ServiceName, ResourceAttributes['service.instance.id'] as ServiceID, Value, TimeUnix 
FROM otel.otel_metrics_gauge 
WHERE TimeUnix >= now() - INTERVAL 1 MINUTE
ORDER BY TimeUnix DESC
FORMAT PrettyCompact
```

```SQL
SELECT
MetricName, ResourceAttributes['service.name'] as ServiceName, ResourceAttributes['service.instance.id'] as ServiceID, Value, TimeUnix 
FROM otel.otel_metrics_gauge 
WHERE MetricName='wave_8' AND ResourceAttributes['service.instance.id']='06444e7bca3b'
ORDER BY TimeUnix DESC
FORMAT PrettyCompact
```