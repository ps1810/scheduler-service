package observer

import (
	"context"
	"scheduler/internal/db/model"
)

type JobEventListener interface {
	OnJobStarted(ctx context.Context, job model.CronJob)
	OnJobCompleted(ctx context.Context, job model.CronJob, result any)
	OnJobFailed(ctx context.Context, job model.CronJob, err error)
}

type JobEventDispatcher struct {
	listeners []JobEventListener
}

func NewJobEventDispatcher() *JobEventDispatcher {
	return &JobEventDispatcher{}
}

func (d *JobEventDispatcher) RegisterListener(listener JobEventListener) {
	d.listeners = append(d.listeners, listener)
}

func (d *JobEventDispatcher) EmitStarted(ctx context.Context, job model.CronJob) {
	for _, l := range d.listeners {
		l.OnJobStarted(ctx, job)
	}
}

func (d *JobEventDispatcher) EmitCompleted(ctx context.Context, job model.CronJob, result any) {
	for _, l := range d.listeners {
		l.OnJobCompleted(ctx, job, result)
	}
}

func (d *JobEventDispatcher) EmitFailed(ctx context.Context, job model.CronJob, err error) {
	for _, l := range d.listeners {
		l.OnJobFailed(ctx, job, err)
	}
}
