package scheduler_strategy

import (
	"fmt"
	"scheduler/internal/db/model"
)

type JobStrategy interface {
	GenerateQuery(job model.CronJob) (string, error)
}

func GetStrategy(duration string) (JobStrategy, error) {
	switch duration {
	case "daily":
		return &DailyStrategy{}, nil
	case "recent_week":
		return &RecentWeekStrategy{}, nil
	case "today", "yesterday", "last_7_days", "last_30_days":
		return &GenericDurationStrategy{}, nil
	default:
		return nil, fmt.Errorf("no strategy for the duration: %s", duration)
	}
}
