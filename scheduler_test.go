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

func TestScheduler_checkLimit(t *testing.T) {
	type fields struct {
		now   time.Time
		limit *Limit
	}
	now, _ := time.Parse("2006-01-02 15:04:05", "2022-10-05 15:30:01")
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "DaysOfWeek in time",
			fields: fields{
				now: now,
				limit: &Limit{
					DaysOfWeek: []time.Weekday{time.Wednesday},
					StartTime:  "00:00",
					EndTime:    "23:59",
					IsBetween:  true,
					When:       nil,
				},
			},
			want: true,
		},
		{
			name: "DaysOfWeek out time",
			fields: fields{
				now: now,
				limit: &Limit{
					DaysOfWeek: []time.Weekday{time.Wednesday},
					StartTime:  "00:00",
					EndTime:    "02:59",
					IsBetween:  true,
					When:       nil,
				},
			},
			want: false,
		},
		{
			name: "DaysOfWeek out time",
			fields: fields{
				now: now,
				limit: &Limit{
					DaysOfWeek: []time.Weekday{time.Wednesday},
					StartTime:  "00:00",
					EndTime:    "02:59",
					IsBetween:  false,
					When:       nil,
				},
			},
			want: true,
		},
		{
			name: "DaysOfWeek out of time",
			fields: fields{
				now: now,
				limit: &Limit{
					DaysOfWeek: []time.Weekday{time.Sunday},
					StartTime:  "00:00",
					EndTime:    "02:59",
					IsBetween:  false,
					When:       nil,
				},
			},
			want: false,
		},
		{
			name: "in time limit",
			fields: fields{
				now: now,
				limit: &Limit{
					DaysOfWeek: []time.Weekday{},
					StartTime:  "00:00",
					EndTime:    "23:59",
					IsBetween:  true,
					When:       nil,
				},
			},
			want: true,
		},
		{
			name: "out of day",
			fields: fields{
				now: now,
				limit: &Limit{
					DaysOfWeek: []time.Weekday{time.Friday},
					StartTime:  "00:00",
					EndTime:    "23:59",
					IsBetween:  true,
					When:       nil,
				},
			},
			want: false,
		},
		{
			name: "in day in time",
			fields: fields{
				now: now,
				limit: &Limit{
					DaysOfWeek: []time.Weekday{time.Friday, time.Monday, time.Wednesday},
					StartTime:  "00:00",
					EndTime:    "23:59",
					IsBetween:  true,
					When:       nil,
				},
			},
			want: true,
		},
		{
			name: "in time limit",
			fields: fields{
				now: now,
				limit: &Limit{
					DaysOfWeek: []time.Weekday{},
					StartTime:  "00:00",
					EndTime:    "23:59",
					IsBetween:  true,
					When: func(ctx context.Context) bool {
						return false
					},
				},
			},
			want: false,
		},
		{
			name: "in time limit",
			fields: fields{
				now: now,
				limit: &Limit{
					DaysOfWeek: []time.Weekday{},
					StartTime:  "00:00",
					EndTime:    "23:59",
					IsBetween:  true,
					When: func(ctx context.Context) bool {
						return true
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scheduler{
				now:   tt.fields.now,
				limit: tt.fields.limit,
			}
			assert.Equalf(t, tt.want, s.checkLimit(), "checkLimit()")
		})
	}
}

func TestScheduler_EveryMinute(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	var marked bool
	ch := make(chan bool, 1)
	s.EveryMinute().CallFunc(func(ctx context.Context) {
		marked = true
		ch <- true
	})
	<-ch
	assert.True(t, marked)
}

func TestScheduler_EveryTwoMinutes(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:30:01")
	s.EveryTwoMinutes()
	assert.Equal(t, s.Next.Minute, s.now.Minute())
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:31:01")
	s.EveryTwoMinutes()
	assert.NotEqual(t, s.Next.Minute, s.now.Minute())
}

func TestScheduler_EveryThreeMinutes(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:30:01")
	s.EveryThreeMinutes()
	assert.Equal(t, s.Next.Minute, s.now.Minute())
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:31:01")
	s.EveryThreeMinutes()
	assert.NotEqual(t, s.Next.Minute, s.now.Minute())
}

func TestScheduler_EveryFourMinutes(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:28:01")
	s.EveryFourMinutes()
	assert.Equal(t, s.Next.Minute, s.now.Minute())
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:31:01")
	s.EveryFourMinutes()
	assert.NotEqual(t, s.Next.Minute, s.now.Minute())
}

func TestScheduler_EveryFiveMinutes(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:30:01")
	s.EveryFiveMinutes()
	assert.Equal(t, s.Next.Minute, s.now.Minute())
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:31:01")
	s.EveryFiveMinutes()
	assert.NotEqual(t, s.Next.Minute, s.now.Minute())
}

func TestScheduler_EveryTenMinutes(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:30:01")
	s.EveryTenMinutes()
	assert.Equal(t, s.Next.Minute, s.now.Minute())
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:31:01")
	s.EveryTenMinutes()
	assert.NotEqual(t, s.Next.Minute, s.now.Minute())
}

func TestScheduler_EveryFifteenMinutes(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:30:01")
	s.EveryFifteenMinutes()
	assert.Equal(t, s.Next.Minute, s.now.Minute())
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:31:01")
	s.EveryFifteenMinutes()
	assert.NotEqual(t, s.Next.Minute, s.now.Minute())
}

func TestScheduler_EveryThirtyMinutes(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:30:01")
	s.EveryThirtyMinutes()
	assert.Equal(t, s.Next.Minute, s.now.Minute())
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:31:01")
	s.EveryThirtyMinutes()
	assert.NotEqual(t, s.Next.Minute, s.now.Minute())
}

func TestScheduler_Hourly(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:00:00")
	s.Hourly()
	assert.Equal(t, s.Next.Minute, 0)
	assert.Equal(t, s.Next.Hour, 15)
}

func TestScheduler_HourlyAt(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:00:00")
	s.HourlyAt(0, 1)
	assert.Equal(t, s.Next.Minute, 0)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:03:00")
	s.HourlyAt(0, 1, 2, 3)
	assert.Equal(t, s.Next.Hour, 15)
	assert.Equal(t, s.Next.Minute, 3)
}
