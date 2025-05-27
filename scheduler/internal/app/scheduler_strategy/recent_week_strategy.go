package scheduler_strategy

import (
	"fmt"
	"scheduler/internal/db/model"
)

type RecentWeekStrategy struct{}

func (r *RecentWeekStrategy) GenerateQuery(job model.CronJob) (string, error) {
	query := fmt.Sprintf(`
WITH week_average AS (
	SELECT %[3]s(%[4]s) AS result, strftime('%%W', %[1]s) as week_number, strftime("%%Y", %[1]s) as year
	FROM %[2]s group by 2,3
)
select result, week_number from week_average where week_number='50' and year='2022'`,
		job.DurationFilter, job.Table, job.Aggregation, job.Field)
	return query, nil
}
