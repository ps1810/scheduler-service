package observer

import (
	"context"
	"go.uber.org/zap"
	"scheduler/internal/db/model"
	"scheduler/internal/logger"
)

type LoggingListener struct{}

func (l *LoggingListener) OnJobStarted(ctx context.Context, job model.CronJob) {
	logger.Log.Info("Job started", zap.String("job_name", job.Name))
}

func (l *LoggingListener) OnJobCompleted(ctx context.Context, job model.CronJob, result any) {
	logger.Log.Info("Job completed", zap.String("job_name", job.Name))
}

func (l *LoggingListener) OnJobFailed(ctx context.Context, job model.CronJob, err error) {
	logger.Log.Error("Job failed", zap.String("job_name", job.Name), zap.Error(err))
}
