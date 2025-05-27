package scheduler_strategy

import (
	"fmt"
	"scheduler/internal/db/model"
	"scheduler/internal/helper"
)

type GenericDurationStrategy struct{}

func (g *GenericDurationStrategy) GenerateQuery(job model.CronJob) (string, error) {
	durationClause := helper.GenerateDurationClauses(job.DurationFilter)
	whereClause, ok := durationClause[job.Duration]
	if !ok {
		return "", fmt.Errorf("unsupported duration: %s", job.Duration)
	}

	query := fmt.Sprintf("SELECT %s(%s) AS result, timestamp FROM %s WHERE %s",
		job.Aggregation, job.Field, job.Table, whereClause)

	return query, nil
}
