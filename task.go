// Package schedule
package schedule

import "context"

type Task interface {
	Run(ctx context.Context)
}

type TaskFunc func(ctx context.Context)

type DefaultTask struct {
	fn TaskFunc
}

func NewDefaultTask(fn TaskFunc) *DefaultTask {
	return &DefaultTask{fn: fn}
}

func (d *DefaultTask) Run(ctx context.Context) {
	d.fn(ctx)
}
