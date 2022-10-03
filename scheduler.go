// Package schedule
package schedule

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Scheduler struct {
	location *time.Location
	now      time.Time
	wg       sync.WaitGroup
	ctx      context.Context
	count    int32
}

func NewScheduler(ctx context.Context, loc *time.Location) *Scheduler {
	return &Scheduler{
		ctx:      ctx,
		location: loc,
		now:      time.Now().In(loc),
		count:    0,
	}
}

func (s *Scheduler) Timezone(loc *time.Location) *Scheduler {
	s.location = loc
	s.now = s.now.In(loc)
	return s
}

func (s *Scheduler) Call(t Task) {
	defer s.Timezone(time.UTC)
	if !s.matchRule() {
		return
	}
	atomic.AddInt32(&s.count, 1)
	s.wg.Add(1)
	go func() {
		defer func() {
			s.wg.Done()
			atomic.AddInt32(&s.count, -1)
			if r := recover(); r != nil {
				log.Println("Recovering schedule task from panic:", r)
			}
		}()
		t.Run(s.ctx)
	}()
}

func (s *Scheduler) CallFunc(fn TaskFunc) {
	s.Call(NewDefaultTask(fn))
}

func (s *Scheduler) matchRule() bool {
	return true
}
