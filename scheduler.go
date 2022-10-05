// Package schedule
package schedule

import (
	"context"
	"github.com/golang-module/carbon/v2"
	"log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Scheduler struct {
	location *time.Location
	now      time.Time
	wg       sync.WaitGroup
	ctx      context.Context
	Next     *NextTick
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
	if !s.isTimeMatched() {
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

func (s *Scheduler) isTimeMatched() bool {
	if s.Next.Omit {
		return false
	}
	if s.Next.Year == s.now.Year() &&
		s.Next.Month == int(s.now.Month()) &&
		s.Next.Day == s.now.Day() &&
		s.Next.Hour == s.now.Hour() &&
		s.Next.Minute == s.now.Minute() {
		return true
	}
	return false
}

func (s *Scheduler) initNextTick() {
	s.Next = &NextTick{
		Year:   s.now.Year(),
		Month:  int(s.now.Month()),
		Day:    s.now.Day(),
		Hour:   s.now.Hour(),
		Minute: 0,
	}
}

func (s *Scheduler) EveryMinute() *Scheduler {
	s.initNextTick()
	s.Next.Minute = s.now.Minute()
	return s
}

func (s *Scheduler) EveryTwoMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%2 == 0 {
		s.Next.Minute = minute
	}
	return s
}

func (s *Scheduler) EveryThreeMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%3 == 0 {
		s.Next.Minute = minute
	}
	return s
}

func (s *Scheduler) EveryFourMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%4 == 0 {
		s.Next.Minute = minute
	}
	return s
}

func (s *Scheduler) EveryFiveMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%5 == 0 {
		s.Next.Minute = minute
	}
	return s
}

func (s *Scheduler) EveryTenMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%10 == 0 {
		s.Next.Minute = minute
	}
	return s
}

func (s *Scheduler) EveryFifteenMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%15 == 0 {
		s.Next.Minute = minute
	}
	return s
}

func (s *Scheduler) EveryThirtyMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%30 == 0 {
		s.Next.Minute = minute
	}
	return s
}

func (s *Scheduler) Hourly() *Scheduler {
	s.initNextTick()
	return s
}

func (s *Scheduler) HourlyAt(t ...int) *Scheduler {
	s.initNextTick()
	s.Next.Omit = true
	minute := s.now.Minute()
	for _, v := range t {
		if v >= 0 && v == minute {
			s.Next.Minute = v
			s.Next.Omit = false
			break
		}
	}
	return s
}

func (s *Scheduler) EveryOddHour() *Scheduler {
	s.initNextTick()
	s.Next.Omit = true
	hour := s.now.Hour()
	if hour >= 1 && hour <= 23 && hour%2 == 0 {
		s.Next.Hour = hour
		s.Next.Omit = false
	}
	return s
}

func (s *Scheduler) setHourlyInterval(n int) {
	s.Next.Omit = true
	hour := s.now.Hour()
	if hour%n == 0 {
		s.Next.Hour = hour
		s.Next.Omit = false
	}
}

func (s *Scheduler) EveryTwoHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(2)
	return s
}

func (s *Scheduler) EveryThreeHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(3)
	return s
}

func (s *Scheduler) EveryFourHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(4)
	return s
}

func (s *Scheduler) EveryFiveHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(5)
	return s
}

func (s *Scheduler) EverySixHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(6)
	return s
}

func (s *Scheduler) Daily() *Scheduler {
	s.initNextTick()
	s.Next.Hour = 0
	return s
}

func (s *Scheduler) setNextTime(t []string) {
	currentHour := s.now.Hour()
	currentMinute := s.now.Minute()
	var hour, minute int
	var err error
	for _, v := range t {
		hm := strings.Split(v, ":")
		if len(hm) == 2 {
			hour, err = strconv.Atoi(hm[0])
			if err == nil {
				minute, err = strconv.Atoi(hm[1])
				if err == nil {
					if currentHour == hour && currentMinute == minute {
						s.Next.Hour = currentHour
						s.Next.Minute = currentMinute
						s.Next.Omit = false
						break
					}
				}
			}
		}
	}
}

func (s *Scheduler) DailyAt(t ...string) *Scheduler {
	s.initNextTick()
	s.Next.Hour = 0
	s.Next.Minute = 0
	s.Next.Omit = true
	s.setNextTime(t)
	return s
}

func (s *Scheduler) TwiceDaily(t ...int) *Scheduler {
	timeList := make([]string, 0, len(t))
	for _, h := range t {
		timeList = append(timeList, strconv.Itoa(h)+":00")
	}
	s.DailyAt(timeList...)
	return s
}

func (s *Scheduler) Weekly() *Scheduler {
	now := carbon.Time2Carbon(s.now)
	now = now.StartOfWeek()
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  now.Month(),
		Day:    now.Day(),
		Hour:   0,
		Minute: 0,
	}
	return s
}

func (s *Scheduler) WeeklyOn(d int, t string) *Scheduler {
	now := carbon.Time2Carbon(s.now)
	now = now.StartOfWeek()
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  now.Month(),
		Day:    0,
		Hour:   0,
		Minute: 0,
		Omit:   true,
	}
	if now.DayOfWeek() == d {
		s.Next.Day = now.Day()
		s.setNextTime([]string{t})
	}
	return s
}

func (s *Scheduler) Monthly() *Scheduler {
	now := carbon.Time2Carbon(s.now)
	now = now.StartOfMonth()
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  now.Month(),
		Day:    now.Day(),
		Hour:   0,
		Minute: 0,
	}
	return s
}

func (s *Scheduler) MonthlyOn(d int, t string) *Scheduler {
	now := carbon.Time2Carbon(s.now)
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  now.Month(),
		Day:    0,
		Hour:   0,
		Minute: 0,
		Omit:   true,
	}
	if now.DayOfMonth() == d {
		s.Next.Day = now.Day()
		s.setNextTime([]string{t})
	}
	return s
}

func (s *Scheduler) TwiceMonthly(b, e int, t string) *Scheduler {
	now := carbon.Time2Carbon(s.now)
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  now.Month(),
		Day:    0,
		Hour:   0,
		Minute: 0,
		Omit:   true,
	}
	day := now.DayOfMonth()
	if day == b || day == e {
		s.Next.Day = day
		s.setNextTime([]string{t})
	}
	return s
}

func (s *Scheduler) LastDayOfMonth(t string) *Scheduler {
	now := carbon.Time2Carbon(s.now)
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  now.Month(),
		Day:    now.EndOfMonth().Day(),
		Hour:   0,
		Minute: 0,
		Omit:   true,
	}
	if t != "" {
		s.setNextTime([]string{t})
	}
	return s
}

func (s *Scheduler) Quarterly() *Scheduler {
	now := carbon.Time2Carbon(s.now)
	qs := now.StartOfQuarter()
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  qs.Month(),
		Day:    qs.Day(),
		Hour:   0,
		Minute: 0,
	}
	return s
}

func (s *Scheduler) Yearly() *Scheduler {
	now := carbon.Time2Carbon(s.now)
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  1,
		Day:    1,
		Hour:   0,
		Minute: 0,
	}
	return s
}

func (s *Scheduler) YearlyOn(m, d int, t string) *Scheduler {
	now := carbon.Time2Carbon(s.now)
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  0,
		Day:    0,
		Hour:   0,
		Minute: 0,
	}
	month := now.Month()
	day := now.Day()
	if month == m && day == d {
		s.Next.Month = month
		s.Next.Day = d
	}
	if t != "" {
		s.setNextTime([]string{t})
	}
	return s
}
