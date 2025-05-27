package scheduler_strategy

import (
	"fmt"
	"scheduler/internal/db/model"
)

type DailyStrategy struct{}

func (d *DailyStrategy) GenerateQuery(job model.CronJob) (string, error) {
	query := fmt.Sprintf("SELECT %s(%s) as total_registration, date(%s) as registration_date FROM %s group by 2",
		job.Aggregation, job.Field, job.DurationFilter, job.Table)
	return query, nil
}
