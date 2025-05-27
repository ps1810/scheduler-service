package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"scheduler/config"
	"scheduler/internal/app/scheduler_strategy"
	"scheduler/internal/db/model"
	"scheduler/internal/helper"
	"scheduler/internal/local_cron"
	"scheduler/internal/logger"
	"scheduler/internal/observer"
	"scheduler/internal/repository"
	"scheduler/pkg/exception"
	"scheduler/pkg/transport"
	"time"
)

// AppScheduler The interface will have the function declaration
type AppScheduler interface {
	LoadAndScheduleJobs(context.Context) error
	AddJob(ctx context.Context, job model.CronJob) error
	RemoveJob(uint) error
	ExecuteJob(context.Context, model.CronJob)
}

// structure to hold the injected object
type appScheduler struct {
	repo       *repository.Repository
	cron       *local_cron.LocalCron
	dispatcher *observer.JobEventDispatcher
}

// NewAppScheduler Inject the objects in the structure
func NewAppScheduler(repo *repository.Repository) AppScheduler {
	d := observer.NewJobEventDispatcher()
	d.RegisterListener(&observer.LoggingListener{})
	return &appScheduler{repo: repo, cron: local_cron.GetCron(), dispatcher: d}
}

func (a *appScheduler) Boot(ctx context.Context) error {
	logger.Log.Info("Booting AppScheduler: loading scheduled jobs")
	return a.LoadAndScheduleJobs(ctx)
}

// LoadAndScheduleJobs Fetching the jobs from database and adding it to the scheduler when service is up
func (a *appScheduler) LoadAndScheduleJobs(ctx context.Context) error {
	jobs, err := a.repo.Job.GetAllEnabledJobs(ctx)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		err := a.AddJob(ctx, job)
		if err != nil {
			logger.Log.Error("Error adding the job to scheduler", zap.String("job_name", job.Name))
		}
	}
	return nil
}

// AddJob Function to add the job to the scheduler. When a new job is created via api it is added to the cron
func (a *appScheduler) AddJob(ctx context.Context, job model.CronJob) error {
	entryID, err := a.cron.Cron.AddFunc(job.CronExpression, func() {
		a.ExecuteJob(ctx, job)
	})
	if err != nil {
		return err
	}
	a.cron.JobMap[job.ID] = entryID
	err = a.UpdateNextRun(job)
	if err != nil {
		logger.Log.Error("Unable to update the next run for the job", zap.String("job_name", job.Name))
	}
	return nil
}

// RemoveJob Removing the job from the scheduler
func (a *appScheduler) RemoveJob(jobID uint) error {
	entryID, ok := a.cron.JobMap[jobID]
	if !ok {
		return exception.NotScheduledJob
	}
	a.cron.Cron.Remove(entryID)
	delete(a.cron.JobMap, jobID)
	return nil
}

// UpdateNextRun Updating the next run of a job
func (a *appScheduler) UpdateNextRun(job model.CronJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	nextRun, err := helper.GetNextRun(job.CronExpression)
	if err != nil {
		logger.Log.Error("Unable to get the next run", zap.String("job_name", job.Name))
	}
	updates := map[string]interface{}{
		"next_run": nextRun.Format("2006-01-02 15:04:05"),
	}
	err = a.repo.Job.UpdateJob(ctx, job, updates)
	if err != nil {
		return err
	}
	return nil
}

// UpdateLastRun updating the last run of the job
func (a *appScheduler) UpdateLastRun(job model.CronJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	updates := map[string]interface{}{
		"last_run": time.Now().Format("2006-01-02 15:04:05"),
	}
	err := a.repo.Job.UpdateJob(ctx, job, updates)
	if err != nil {
		return err
	}
	return nil
}

// ExecuteJob Function to execute the job. Currently part of the scheduler app. It can be moved out if different kind of jobs are to be executed
func (a *appScheduler) ExecuteJob(ctx context.Context, job model.CronJob) {
	//logger.Log.Info("Executing job", zap.String("job_name", job.Name))
	a.dispatcher.EmitStarted(ctx, job)
	err := a.UpdateLastRun(job)
	err = a.UpdateNextRun(job)
	if err != nil {
		logger.Log.Error("Unable to update the last run for the job", zap.String("job_name", job.Name))
	}
	strategyImpl, err := scheduler_strategy.GetStrategy(job.Duration)
	if err != nil {
		logger.Log.Error("Duration of Job is unknown")
		return
	}

	query, err := strategyImpl.GenerateQuery(job)
	if err != nil {
		logger.Log.Error("Error generating query")
	}
	data := a.repo.Job.ExecuteRawQuery(query)
	if data == nil {
		a.dispatcher.EmitCompleted(ctx, job, nil)
		//logger.Log.Info("No data found for the query")
	} else {
		result := map[string]interface{}{
			"result": data,
		}
		err = PostResult(result)
		if err != nil {
			a.dispatcher.EmitFailed(ctx, job, err)
			//logger.Log.Error("Unable to send the data")
		}
	}
}

// PostResult Function to call the api to post result
func PostResult(data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	requestConfig := config.GetConfig().PostResult
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req := transport.HttpRequest{
		HttpClient: transport.NewHTTPClient(),
		Method:     requestConfig.Method,
		Url:        fmt.Sprintf("%s:%d/%s", requestConfig.Url, requestConfig.Port, requestConfig.Path),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	}
	err = transport.RequestAndParseJSONBody(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
