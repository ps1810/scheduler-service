package jobs

import (
	"context"
	"scheduler/internal/app/scheduler"
	"scheduler/internal/db/model"
	"scheduler/internal/helper"
	"scheduler/internal/interface/request"
	"scheduler/internal/interface/response"
	"scheduler/internal/repository"
	"scheduler/pkg/exception"
	"time"
)

type JobsApp interface {
	AddJob(context.Context, request.SchedulerRequest) error
	GetAllJobs(context.Context) ([]response.Jobs, error)
	DeleteAJob(context.Context, int) error
}

type jobsAppImpl struct {
	Repo *repository.Repository
	sch  scheduler.AppScheduler
}

// NewJobApp Injecting the repo and scheduler object
func NewJobApp(repo *repository.Repository) JobsApp {
	return &jobsAppImpl{Repo: repo, sch: scheduler.NewAppScheduler(repo)}
}

// GetAllJobs returns all the jobs from the database
func (j *jobsAppImpl) GetAllJobs(ctx context.Context) ([]response.Jobs, error) {
	jobs, err := j.Repo.Job.GetAllJobs(ctx)
	if err != nil {
		return nil, err
	}
	var resp []response.Jobs
	for _, job := range jobs {
		resp = append(resp, response.Jobs{
			ID:             job.ID,
			Name:           job.Name,
			CronExpression: job.CronExpression,
			Enabled:        job.Enabled,
			Table:          job.Table,
			Field:          job.Field,
			Aggregation:    job.Aggregation,
			Duration:       job.Duration,
			DurationFilter: job.DurationFilter,
			CreatedAt:      job.CreatedAt,
			LastRun:        job.LastRun,
			NextRun:        job.NextRun,
		})
	}
	return resp, nil
}

// DeleteAJob delete a job from a database. It basically set the enabled as false
func (j *jobsAppImpl) DeleteAJob(ctx context.Context, jobID int) error {
	job, err := j.Repo.Job.GetAJobFromID(ctx, jobID)
	if err != nil {
		return exception.DataNotFoundError
	}
	updates := map[string]interface{}{
		"enabled": false,
	}
	err = j.Repo.Job.UpdateJob(ctx, job, updates)
	if err != nil {
		return exception.UpdateFailedError
	}
	err = j.sch.RemoveJob(job.ID)
	if err != nil {
		return err
	}
	return nil
}

// AddJob Add a job in the database and to the scheduler for execution
func (j *jobsAppImpl) AddJob(ctx context.Context, requestBody request.SchedulerRequest) error {
	nextRun, err := helper.GetNextRun(requestBody.CronSchedule)
	if err != nil {
		return err
	}
	newJob := model.CronJob{
		Name:           requestBody.Name,
		CronExpression: requestBody.CronSchedule,
		Enabled:        true,
		Table:          requestBody.Table,
		Field:          requestBody.Field,
		Aggregation:    requestBody.Aggregation,
		Duration:       requestBody.DurationOption,
		DurationFilter: requestBody.DurationFilter,
		CreatedAt:      time.Now().UTC().Format("2006-01-02 15:04:05"),
		LastRun:        "",
		NextRun:        nextRun.Format("2006-01-02 15:04:05"),
	}

	err = j.Repo.Job.AddAJob(ctx, &newJob)
	if err != nil {
		return exception.FailedAddJobError
	}
	err = j.sch.AddJob(ctx, newJob)
	if err != nil {
		return err
	}
	return nil
}
