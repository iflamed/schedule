// Package schedule
package schedule

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScheduler_Timezone(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	prc, err := time.LoadLocation("Asia/Shanghai")
	assert.NoError(t, err)
	hour := s.now.Hour()
	s.Timezone(prc)
	assert.NotEqual(t, hour, s.now.Hour())
}

func TestScheduler_isTimeMatched(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Next = &NextTick{
		Year:   s.now.Year(),
		Month:  int(s.now.Month()),
		Day:    s.now.Day(),
		Hour:   s.now.Hour(),
		Minute: s.now.Minute(),
		Omit:   true,
	}
	assert.False(t, s.isTimeMatched())
	s.Next.Omit = false
	assert.True(t, s.isTimeMatched())
	s.Next.Minute = 60
	assert.False(t, s.isTimeMatched())
}

func TestScheduler_timeToMinutes(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	var hour, minute int
	hour, minute = s.timeToMinutes("a:b")
	assert.Zero(t, hour)
	assert.Zero(t, minute)
	hour, minute = s.timeToMinutes("a:1")
	assert.Zero(t, hour)
	assert.Zero(t, minute)
	hour, minute = s.timeToMinutes("1:b")
	assert.Zero(t, hour)
	assert.Zero(t, minute)
	hour, minute = s.timeToMinutes("1:1")
	assert.Equal(t, 1, hour)
	assert.Equal(t, 1, minute)
}
