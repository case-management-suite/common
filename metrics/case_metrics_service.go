package metrics

import (
	"github.com/rs/zerolog/log"
)

type caseMetricsType struct {
	CaseCreated       MetricEvent
	CaseClosed        MetricEvent
	CaseStatusChanged MetricEvent
	UnsupportedAction MetricEvent
}

var CaseMetrics = caseMetricsType{
	CaseStatusChanged: "CASE_STATUS_CHANGED",
	CaseCreated:       "CASE_CREATED",
	CaseClosed:        "CASE_CLOSED",
	UnsupportedAction: "ACTION_NOT_SUPPORTED",
}

type CaseMetricsService struct{}

func NewCaseMetricsService() MetricsService {
	return CaseMetricsService{}
}

func (cms CaseMetricsService) LogEvent(event MetricEvent) {
	log.Logger.Print(event)

}
