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

func TestScheduler_EveryOddHour(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:00:00")
	s.EveryOddHour()
	assert.Equal(t, s.Next.Hour, 15)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 00:00:00")
	s.EveryOddHour()
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 02:00:00")
	s.EveryOddHour()
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
}

func TestScheduler_EveryTwoHours(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 15:00:00")
	s.EveryTwoHours()
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 00:00:00")
	s.EveryTwoHours()
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 02:00:00")
	s.EveryTwoHours()
	assert.Equal(t, s.Next.Hour, 2)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_EveryThreeHours(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 14:00:00")
	s.EveryThreeHours()
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 00:00:00")
	s.EveryThreeHours()
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 03:00:00")
	s.EveryThreeHours()
	assert.Equal(t, s.Next.Hour, 3)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_EveryFourHours(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 14:00:00")
	s.EveryFourHours()
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 00:00:00")
	s.EveryFourHours()
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 04:00:00")
	s.EveryFourHours()
	assert.Equal(t, s.Next.Hour, 4)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_EveryFiveHours(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 16:00:00")
	s.EveryFiveHours()
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 00:00:00")
	s.EveryFiveHours()
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 05:00:00")
	s.EveryFiveHours()
	assert.Equal(t, s.Next.Hour, 5)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_EverySixHours(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 17:00:00")
	s.EverySixHours()
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 00:00:00")
	s.EverySixHours()
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 06:00:00")
	s.EverySixHours()
	assert.Equal(t, s.Next.Hour, 6)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_Daily(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 17:00:00")
	s.Daily()
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_DailyAt(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 17:01:00")
	s.DailyAt("00:30", "17:01")
	assert.Equal(t, s.Next.Hour, 17)
	assert.Equal(t, s.Next.Minute, 1)
	assert.False(t, s.Next.Omit)
	s.DailyAt("00:30", "16:01")
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
}

func TestScheduler_At(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 17:01:00")
	s.At("00:30", "17:01")
	assert.Equal(t, s.Next.Hour, 17)
	assert.Equal(t, s.Next.Minute, 1)
	assert.False(t, s.Next.Omit)
	s.At("00:30", "16:01")
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
}

func TestScheduler_TwiceDaily(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 17:00:00")
	s.TwiceDaily(1, 17)
	assert.Equal(t, s.Next.Hour, 17)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 01:00:00")
	s.TwiceDaily(1, 17)
	assert.Equal(t, s.Next.Hour, 1)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 03:00:00")
	s.TwiceDaily(1, 17)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
}

func TestScheduler_TwiceDailyAt(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 17:05:00")
	s.TwiceDailyAt(1, 17, 5)
	assert.Equal(t, s.Next.Hour, 17)
	assert.Equal(t, s.Next.Minute, 5)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 01:04:00")
	s.TwiceDailyAt(1, 17, 4)
	assert.Equal(t, s.Next.Hour, 1)
	assert.Equal(t, s.Next.Minute, 4)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 03:02:00")
	s.TwiceDailyAt(1, 17, 2)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
}

func TestScheduler_Weekly(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-05 17:00:00")
	s.Weekly()
	assert.Equal(t, s.Next.Day, 2)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_WeeklyOn(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-04 17:00:00")
	s.WeeklyOn(time.Tuesday, "16:10")
	assert.Equal(t, s.Next.Day, 4)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
	s.WeeklyOn(time.Tuesday, "17:00")
	assert.Equal(t, s.Next.Day, 4)
	assert.Equal(t, s.Next.Hour, 17)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_Monthly(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-01 00:00:00")
	s.Monthly()
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 1)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)

	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-02 00:00:00")
	s.Monthly()
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 1)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_MonthlyOn(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-02 03:00:00")
	s.MonthlyOn(2, "03:00")
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 2)
	assert.Equal(t, s.Next.Hour, 3)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.MonthlyOn(3, "00:00")
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 0)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
}

func TestScheduler_TwiceMonthly(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-02 03:00:00")
	s.TwiceMonthly(2, 3, "03:00")
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 2)
	assert.Equal(t, s.Next.Hour, 3)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-03 03:00:00")
	s.TwiceMonthly(2, 3, "03:00")
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 3)
	assert.Equal(t, s.Next.Hour, 3)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-03 04:00:00")
	s.TwiceMonthly(2, 3, "03:00")
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 3)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
}

func TestScheduler_LastDayOfMonth(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-02 03:00:00")
	s.LastDayOfMonth("03:00")
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 31)
	assert.Equal(t, s.Next.Hour, 3)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)

	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-02 03:00:00")
	s.LastDayOfMonth("02:00")
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 31)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
}

func TestScheduler_Quarterly(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-10 03:00:00")
	s.Quarterly()
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 1)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_Yearly(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-10 03:00:00")
	s.Yearly()
	assert.Equal(t, s.Next.Month, 1)
	assert.Equal(t, s.Next.Day, 1)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
}

func TestScheduler_YearlyOn(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-10 03:00:00")
	s.YearlyOn(10, 10, "03:00")
	assert.Equal(t, s.Next.Month, 10)
	assert.Equal(t, s.Next.Day, 10)
	assert.Equal(t, s.Next.Hour, 3)
	assert.Equal(t, s.Next.Minute, 0)
	assert.False(t, s.Next.Omit)
	s.now, _ = time.Parse("2006-01-02 15:04:05", "2022-10-10 03:00:00")
	s.YearlyOn(8, 9, "02:00")
	assert.Equal(t, s.Next.Month, 0)
	assert.Equal(t, s.Next.Day, 0)
	assert.Equal(t, s.Next.Hour, 0)
	assert.Equal(t, s.Next.Minute, 0)
	assert.True(t, s.Next.Omit)
}

func TestScheduler_Weekdays(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Weekdays()
	assert.Len(t, s.limit.DaysOfWeek, 5)
	assert.Contains(t, s.limit.DaysOfWeek, time.Monday)
	assert.Contains(t, s.limit.DaysOfWeek, time.Tuesday)
	assert.Contains(t, s.limit.DaysOfWeek, time.Wednesday)
	assert.Contains(t, s.limit.DaysOfWeek, time.Thursday)
	assert.Contains(t, s.limit.DaysOfWeek, time.Friday)
	assert.NotContains(t, s.limit.DaysOfWeek, time.Saturday)
	assert.NotContains(t, s.limit.DaysOfWeek, time.Sunday)
}

func TestScheduler_Weekends(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Weekends()
	assert.Len(t, s.limit.DaysOfWeek, 2)
	assert.Contains(t, s.limit.DaysOfWeek, time.Saturday)
	assert.Contains(t, s.limit.DaysOfWeek, time.Sunday)
}

func TestScheduler_Mondays(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Mondays()
	assert.Len(t, s.limit.DaysOfWeek, 1)
	assert.Contains(t, s.limit.DaysOfWeek, time.Monday)
}

func TestScheduler_Tuesdays(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Tuesdays()
	assert.Len(t, s.limit.DaysOfWeek, 1)
	assert.Contains(t, s.limit.DaysOfWeek, time.Tuesday)
}

func TestScheduler_Wednesdays(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Wednesdays()
	assert.Len(t, s.limit.DaysOfWeek, 1)
	assert.Contains(t, s.limit.DaysOfWeek, time.Wednesday)
}

func TestScheduler_Thursdays(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Thursdays()
	assert.Len(t, s.limit.DaysOfWeek, 1)
	assert.Contains(t, s.limit.DaysOfWeek, time.Thursday)
}

func TestScheduler_Fridays(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Fridays()
	assert.Len(t, s.limit.DaysOfWeek, 1)
	assert.Contains(t, s.limit.DaysOfWeek, time.Friday)
}

func TestScheduler_Saturdays(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Saturdays()
	assert.Len(t, s.limit.DaysOfWeek, 1)
	assert.Contains(t, s.limit.DaysOfWeek, time.Saturday)
}

func TestScheduler_Sundays(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Sundays()
	assert.Len(t, s.limit.DaysOfWeek, 1)
	assert.Contains(t, s.limit.DaysOfWeek, time.Sunday)
}

func TestScheduler_Days(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Days(time.Monday, time.Friday)
	assert.Len(t, s.limit.DaysOfWeek, 2)
	assert.Contains(t, s.limit.DaysOfWeek, time.Monday)
	assert.Contains(t, s.limit.DaysOfWeek, time.Friday)
	assert.NotContains(t, s.limit.DaysOfWeek, time.Wednesday)
}

func TestScheduler_Between(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.Between("09:00", "15:00")
	assert.Equal(t, "09:00", s.limit.StartTime)
	assert.Equal(t, "15:00", s.limit.EndTime)
	assert.True(t, s.limit.IsBetween)
}

func TestScheduler_UnlessBetween(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.UnlessBetween("09:00", "15:00")
	assert.Equal(t, "09:00", s.limit.StartTime)
	assert.Equal(t, "15:00", s.limit.EndTime)
	assert.False(t, s.limit.IsBetween)
}

func TestScheduler_When(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	s.When(func(ctx context.Context) bool {
		return false
	})
	assert.False(t, s.limit.When(s.ctx))
	s.When(func(ctx context.Context) bool {
		return true
	})
	assert.True(t, s.limit.When(s.ctx))
}

func TestScheduler_Call(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	dayTime := s.now.Format("15:04")
	ch := make(chan bool, 1)
	s.DailyAt(dayTime).When(func(ctx context.Context) bool {
		return true
	}).Call(NewDefaultTask(func(ctx context.Context) {
		ch <- true
	}))
	assert.True(t, <-ch)
}
