package metrics

type MetricEvent string

type MetricsService interface {
	LogEvent(MetricEvent)
}
