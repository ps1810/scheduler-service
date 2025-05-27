package helper

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"scheduler/pkg/exception"
	"strings"
	"time"
)

func GenerateDurationClauses(columnName string) map[string]string {
	return map[string]string{
		"today":        fmt.Sprintf("date(%s) = date('2022-12-22')", columnName),
		"yesterday":    fmt.Sprintf("date(%s) = date('2022-12-22', '-1 day')", columnName),
		"last_7_days":  fmt.Sprintf("date(%s) = date('2022-12-22', '-7 day')", columnName),
		"last_30_days": fmt.Sprintf("date(%s) = date('2022-12-22', '-1 month')", columnName),
	}
}

func parseCronExpression(expr string) (cron.Schedule, error) {
	// If the expression starts with "@", use descriptor-enabled parser
	if strings.HasPrefix(expr, "@") {
		parser := cron.NewParser(
			cron.Second | cron.Minute | cron.Hour |
				cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		)
		return parser.Parse(expr)
	}

	// Otherwise use standard 5-field parser (no seconds or descriptors)
	return cron.ParseStandard(expr)
}

func GetNextRun(cronExpression string) (time.Time, error) {
	schedule, err := parseCronExpression(cronExpression)
	if err != nil {
		return time.Time{}, exception.InvalidCronExpression
	}
	nextRun := schedule.Next(time.Now().UTC())
	return nextRun, nil
}
