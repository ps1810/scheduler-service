package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"scheduler/internal/db/model"
)

// JobRepository Function declaration for running the query in database
type JobRepository interface {
	GetAllEnabledJobs(context.Context) ([]model.CronJob, error)
	UpdateJob(context.Context, model.CronJob, map[string]interface{}) error
	GetAllJobs(context.Context) ([]model.CronJob, error)
	AddAJob(context.Context, *model.CronJob) error
	GetAJobFromID(context.Context, int) (model.CronJob, error)
	ExecuteRawQuery(string) []map[string]interface{}
}

type JobRepositoryImpl struct {
	DB *gorm.DB
}

// NewJobRepository Function to inject the database object
func NewJobRepository(db *gorm.DB) JobRepository {
	return &JobRepositoryImpl{DB: db}
}

// GetAllEnabledJobs Get all the enbaled jobs from the database
func (j *JobRepositoryImpl) GetAllEnabledJobs(ctx context.Context) ([]model.CronJob, error) {
	var jobs []model.CronJob
	err := j.DB.WithContext(ctx).Where("enabled = ?", true).Find(&jobs).Error
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

// UpdateJob Update field for a job
func (j *JobRepositoryImpl) UpdateJob(ctx context.Context, job model.CronJob, updates map[string]interface{}) error {
	err := j.DB.WithContext(ctx).Model(&model.CronJob{}).Where("id = ?", job.ID).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

// GetAllJobs Get all jobs from the database
func (j *JobRepositoryImpl) GetAllJobs(ctx context.Context) ([]model.CronJob, error) {
	var jobs []model.CronJob
	err := j.DB.WithContext(ctx).Model(&model.CronJob{}).Find(&jobs).Error
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

// AddAJob Add a job in the database
func (j *JobRepositoryImpl) AddAJob(ctx context.Context, job *model.CronJob) error {
	err := j.DB.WithContext(ctx).Create(job).Error
	if err != nil {
		return err
	}
	return nil
}

// GetAJobFromID Get job from id
func (j *JobRepositoryImpl) GetAJobFromID(ctx context.Context, jobID int) (model.CronJob, error) {
	var job model.CronJob
	err := j.DB.WithContext(ctx).Where("id = ?", jobID).First(&job).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return job, err
	}
	if err != nil {
		return job, err
	}
	return job, nil
}

// ExecuteRawQuery Execute raw query directly
func (j *JobRepositoryImpl) ExecuteRawQuery(query string) []map[string]interface{} {
	var rows []map[string]interface{}
	err := j.DB.Raw(query).Scan(&rows).Error
	if err != nil {
		return nil
	}
	return rows
}
