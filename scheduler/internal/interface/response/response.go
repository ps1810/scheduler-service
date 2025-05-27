package response

import (
	"scheduler/internal/interface/request"
	"scheduler/pkg/exception"
)

type CommonResponse struct {
	ResponseCode    int                        `json:"response_code"`
	ResponseMessage string                     `json:"response_message"`
	Errors          *exception.ExceptionErrors `json:"errors,omitempty"`
	Data            any                        `json:"data,omitempty"`
}

type Jobs struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	CronExpression string `json:"cron_expression"`
	Enabled        bool   `json:"enabled"`
	Table          string `json:"table"`
	Field          string `json:"field"`
	Aggregation    string `json:"aggregation"`
	Duration       string `json:"duration"`
	DurationFilter string `json:"duration_filter"`
	CreatedAt      string `json:"created_at"`
	LastRun        string `json:"last_run"`
	NextRun        string `json:"next_run"`
}

type MetricConfig struct {
	Table          string   `json:"table"`
	Fields         []string `json:"fields"`
	Aggregations   []string `json:"aggregations"`
	DurationFilter []string `json:"duration_filter"`
	Durations      []string `json:"durations"`
}

type MetadataResponse struct {
	Metrics   []MetricConfig `json:"metrics"`
	Durations []string       `json:"durations"`
}

var Metrics = []MetricConfig{
	{
		Table:          "registration",
		Fields:         []string{"weight"},
		Aggregations:   []string{"min", "max", "avg", "count"},
		DurationFilter: []string{"timestamp"},
		Durations:      request.DurationOptions,
	},
}
