package request

import (
	"errors"
)

type SchedulerRequest struct {
	Name           string `json:"name"`
	Table          string `json:"table"`
	Field          string `json:"field"`
	Aggregation    string `json:"aggregation"`
	DurationFilter string `json:"duration_filter"`
	DurationOption string `json:"duration_option"`
	CronSchedule   string `json:"cron_schedule"`
}

var AllowedTables = map[string][]string{
	"registration": {"weight"},
}

var AllowedAggregations = map[string]string{
	"min":   "MIN",
	"max":   "MAX",
	"avg":   "AVG",
	"count": "COUNT",
}

var DurationOptions = []string{"today", "yesterday", "last_7_days", "last_30_days", "recent_week", "daily"}

func (sch *SchedulerRequest) ValidateSchedulerRequest() error {
	fields, tableOk := AllowedTables[sch.Table]
	if !tableOk {
		return errors.New("invalid Table")
	}

	fieldOk := false
	for _, f := range fields {
		if f == sch.Field {
			fieldOk = true
			break
		}
	}

	if !fieldOk {
		return errors.New("invalid field for the table")
	}

	_, aggOk := AllowedAggregations[sch.Aggregation]
	if !aggOk {
		return errors.New("invalid aggregation")
	}

	durOk := false
	for _, dur := range DurationOptions {
		if dur == sch.DurationOption {
			durOk = true
			break
		}
	}

	if !durOk {
		return errors.New("invalid duration")
	}

	return nil
}
