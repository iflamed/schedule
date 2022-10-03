// Package schedule
package schedule

import "context"

type Task interface {
	Run(ctx context.Context)
}

type TaskFunc func(ctx context.Context)
