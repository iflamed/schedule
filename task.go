// Package schedule
package schedule

import (
	"context"
	"log"
	"time"
)

type Task interface {
	Run(ctx context.Context)
}

type Logger interface {
	Error(msg string, e any)
	Debugf(msg string, n int32)
	Debug(msg string)
}

type TaskFunc func(ctx context.Context)
type WhenFunc func(ctx context.Context) bool

type DefaultTask struct {
	fn TaskFunc
}

func NewDefaultTask(fn TaskFunc) *DefaultTask {
	return &DefaultTask{fn: fn}
}

func (d *DefaultTask) Run(ctx context.Context) {
	d.fn(ctx)
}

type NextTick struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Omit   bool
}

type Limit struct {
	DaysOfWeek []time.Weekday
	StartTime  string
	EndTime    string
	IsBetween  bool
	When       WhenFunc
}

type DefaultLogger struct {
}

func (d *DefaultLogger) Error(msg string, r any) {
	log.Println(msg, r)
}

func (d *DefaultLogger) Debug(msg string) {
	log.Println(msg)
}

func (d *DefaultLogger) Debugf(msg string, i int32) {
	log.Printf(msg, i)
}
