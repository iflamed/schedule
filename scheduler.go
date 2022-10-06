// Package schedule
// The core code of scheduler, contain frequency options and constraints.
package schedule

import (
	"context"
	"github.com/golang-module/carbon/v2"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Scheduler The core scheduler struct
type Scheduler struct {
	location *time.Location
	now      time.Time
	wg       sync.WaitGroup
	ctx      context.Context
	Next     *NextTick
	limit    *Limit
	count    int32
	log      Logger
}

// NewScheduler create instance of scheduler with context and default time.location
func NewScheduler(ctx context.Context, loc *time.Location) *Scheduler {
	return &Scheduler{
		ctx:      ctx,
		location: loc,
		now:      time.Now().In(loc),
		Next:     &NextTick{},
		limit:    &Limit{},
		count:    0,
		log:      &DefaultLogger{},
	}
}

// Timezone set timezone with a new time.Location instance
// after `Call` and `CallFunc` method called, the current time will roll back to default location.
func (s *Scheduler) Timezone(loc *time.Location) *Scheduler {
	s.now = s.now.In(loc)
	return s
}

// SetLogger set a new logger
func (s *Scheduler) SetLogger(l Logger) *Scheduler {
	if l == nil {
		return s
	}
	s.log = l
	return s
}

// Start wait all task to be finished
func (s *Scheduler) Start() {
	if atomic.LoadInt32(&s.count) > 0 {
		s.log.Debugf("Wait for %d tasks finish... \n", s.count)
	}
	s.wg.Wait()
	s.log.Debug("All tasks have been finished.")
}

// Call call a task
func (s *Scheduler) Call(t Task) {
	defer s.Timezone(s.location)
	if !s.isTimeMatched() {
		return
	}
	if !s.checkLimit() {
		return
	}
	atomic.AddInt32(&s.count, 1)
	s.wg.Add(1)
	go func() {
		defer func() {
			s.wg.Done()
			atomic.AddInt32(&s.count, -1)
			if r := recover(); r != nil {
				s.log.Error("Recovering schedule task from panic:", r)
			}
		}()
		t.Run(s.ctx)
	}()
}

// CallFunc call a task function
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

func (s *Scheduler) timeToMinutes(t string) (hour, minute int) {
	var err error
	hm := strings.Split(t, ":")
	if len(hm) == 2 {
		hour, err = strconv.Atoi(hm[0])
		if err == nil {
			minute, err = strconv.Atoi(hm[1])
		}
	}
	if err != nil {
		hour = 0
		minute = 0
	}
	return
}

func (s *Scheduler) checkLimit() bool {
	if len(s.limit.DaysOfWeek) > 0 {
		var inDays bool
		for _, day := range s.limit.DaysOfWeek {
			if day == s.now.Weekday() {
				inDays = true
			}
		}
		if !inDays {
			return false
		}
	}
	var startMinute, endMinute int
	var hour, minute int
	if s.limit.StartTime != "" {
		hour, minute = s.timeToMinutes(s.limit.StartTime)
		startMinute = hour*60 + minute
	}
	if s.limit.EndTime != "" {
		hour, minute = s.timeToMinutes(s.limit.EndTime)
		endMinute = hour*60 + minute
	}
	if startMinute > endMinute {
		temp := startMinute
		startMinute = endMinute
		endMinute = temp
	}
	minuteOffset := s.now.Hour()*60 + s.now.Minute()
	if s.limit.IsBetween && (minuteOffset < startMinute || minuteOffset > endMinute) {
		return false
	} else if !s.limit.IsBetween && minuteOffset > startMinute && minuteOffset < endMinute {
		return false
	}

	if s.limit.When != nil {
		return s.limit.When(s.ctx)
	}
	return true
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

// EveryMinute run task every minutes
func (s *Scheduler) EveryMinute() *Scheduler {
	s.initNextTick()
	s.Next.Minute = s.now.Minute()
	return s
}

// EveryTwoMinutes run task every two minutes
func (s *Scheduler) EveryTwoMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%2 == 0 {
		s.Next.Minute = minute
	}
	return s
}

// EveryThreeMinutes run task every three minutes
func (s *Scheduler) EveryThreeMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%3 == 0 {
		s.Next.Minute = minute
	}
	return s
}

// EveryFourMinutes run task every four minutes
func (s *Scheduler) EveryFourMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%4 == 0 {
		s.Next.Minute = minute
	}
	return s
}

// EveryFiveMinutes run task every five minutes
func (s *Scheduler) EveryFiveMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%5 == 0 {
		s.Next.Minute = minute
	}
	return s
}

// EveryTenMinutes run the task every ten minutes
func (s *Scheduler) EveryTenMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%10 == 0 {
		s.Next.Minute = minute
	}
	return s
}

// EveryFifteenMinutes run the task every fifteen minutes
func (s *Scheduler) EveryFifteenMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%15 == 0 {
		s.Next.Minute = minute
	}
	return s
}

// EveryThirtyMinutes run the task every thirty minutes
func (s *Scheduler) EveryThirtyMinutes() *Scheduler {
	s.initNextTick()
	minute := s.now.Minute()
	if minute%30 == 0 {
		s.Next.Minute = minute
	}
	return s
}

// Hourly run the task every hour
func (s *Scheduler) Hourly() *Scheduler {
	s.initNextTick()
	return s
}

// HourlyAt run the task every hour at some minutes past the hour
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

// EveryOddHour run the task every odd hour
func (s *Scheduler) EveryOddHour() *Scheduler {
	s.initNextTick()
	s.Next.Omit = true
	hour := s.now.Hour()
	if hour >= 1 && hour <= 23 && hour%2 != 0 {
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

// EveryTwoHours run the task every two hours
func (s *Scheduler) EveryTwoHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(2)
	return s
}

// EveryThreeHours run the task every three hours
func (s *Scheduler) EveryThreeHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(3)
	return s
}

// EveryFourHours run the task every four hours
func (s *Scheduler) EveryFourHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(4)
	return s
}

// EveryFiveHours run the task every five hours
func (s *Scheduler) EveryFiveHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(5)
	return s
}

// EverySixHours run the task every six hours
func (s *Scheduler) EverySixHours() *Scheduler {
	s.initNextTick()
	s.setHourlyInterval(6)
	return s
}

// Daily run the task every day at midnight
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

// At run the task every day at some time (03:00 format), method alias of dailyAt
func (s *Scheduler) At(t ...string) *Scheduler {
	return s.DailyAt(t...)
}

// DailyAt run the task every day at some time (03:00 format)
func (s *Scheduler) DailyAt(t ...string) *Scheduler {
	s.initNextTick()
	s.Next.Hour = 0
	s.Next.Minute = 0
	s.Next.Omit = true
	s.setNextTime(t)
	return s
}

// TwiceDaily run the task daily at first and second hour
func (s *Scheduler) TwiceDaily(first, second int) *Scheduler {
	timeList := make([]string, 0, 2)
	timeList = append(timeList, strconv.Itoa(first)+":00")
	timeList = append(timeList, strconv.Itoa(second)+":00")
	s.DailyAt(timeList...)
	return s
}

// TwiceDailyAt run the task daily at some time
// TwiceDailyAt(1, 13, 15) run the task daily at 1:15 & 13:15
func (s *Scheduler) TwiceDailyAt(first, second, offset int) *Scheduler {
	timeList := make([]string, 0, 2)
	timeList = append(timeList, strconv.Itoa(first)+":"+strconv.Itoa(offset))
	timeList = append(timeList, strconv.Itoa(second)+":"+strconv.Itoa(offset))
	s.DailyAt(timeList...)
	return s
}

// Weekly run the task every Sunday at 00:00
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

// WeeklyOn run the task every week on a time
// WeeklyOn(1, "8:00") run the task every week on Monday at 8:00
func (s *Scheduler) WeeklyOn(d time.Weekday, t string) *Scheduler {
	s.Next = &NextTick{
		Year:   s.now.Year(),
		Month:  int(s.now.Month()),
		Day:    0,
		Hour:   0,
		Minute: 0,
		Omit:   true,
	}
	if s.now.Weekday() == d {
		s.Next.Day = s.now.Day()
		s.setNextTime([]string{t})
	}
	return s
}

// Monthly run the task on the first day of every month at 00:00
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

// MonthlyOn run the task every month on a time
// MonthlyOn(4, "15:00") run the task every month on the 4th at 15:00
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

// TwiceMonthly run the task monthly on some time
// TwiceMonthly(1, 16, "13:00") run the task monthly on the 1st and 16th at 13:00
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

// LastDayOfMonth run the task on the last day of the month at a time
// LastDayOfMonth("15:00") run the task on the last day of the month at 15:00
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

// Quarterly Run the task on the first day of every quarter at 00:00
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

// Yearly run the task on the first day of every year at 00:00
func (s *Scheduler) Yearly() *Scheduler {
	s.Next = &NextTick{
		Year:   s.now.Year(),
		Month:  1,
		Day:    1,
		Hour:   0,
		Minute: 0,
	}
	return s
}

// YearlyOn Run the task every year on a time
// YearlyOn(6, 1, "17:00") run the task every year on June 1st at 17:00
func (s *Scheduler) YearlyOn(m, d int, t string) *Scheduler {
	now := carbon.Time2Carbon(s.now)
	s.Next = &NextTick{
		Year:   now.Year(),
		Month:  0,
		Day:    0,
		Hour:   0,
		Minute: 0,
		Omit:   true,
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

// Weekdays limit the task to weekdays
func (s *Scheduler) Weekdays() *Scheduler {
	s.limit.DaysOfWeek = append(
		s.limit.DaysOfWeek,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
	)
	return s
}

// Weekends limit the task to weekends
func (s *Scheduler) Weekends() *Scheduler {
	s.limit.DaysOfWeek = append(
		s.limit.DaysOfWeek,
		time.Saturday,
		time.Sunday,
	)
	return s
}

// Mondays limit the task to Monday
func (s *Scheduler) Mondays() *Scheduler {
	s.limit.DaysOfWeek = append(
		s.limit.DaysOfWeek,
		time.Monday,
	)
	return s
}

// Tuesdays limit the task to Tuesday
func (s *Scheduler) Tuesdays() *Scheduler {
	s.limit.DaysOfWeek = append(
		s.limit.DaysOfWeek,
		time.Tuesday,
	)
	return s
}

// Wednesdays limit the task to Wednesday
func (s *Scheduler) Wednesdays() *Scheduler {
	s.limit.DaysOfWeek = append(
		s.limit.DaysOfWeek,
		time.Wednesday,
	)
	return s
}

// Thursdays limit the task to Thursday
func (s *Scheduler) Thursdays() *Scheduler {
	s.limit.DaysOfWeek = append(
		s.limit.DaysOfWeek,
		time.Thursday,
	)
	return s
}

// Fridays limit the task to Friday
func (s *Scheduler) Fridays() *Scheduler {
	s.limit.DaysOfWeek = append(
		s.limit.DaysOfWeek,
		time.Friday,
	)
	return s
}

// Saturdays limit the task to Saturday
func (s *Scheduler) Saturdays() *Scheduler {
	s.limit.DaysOfWeek = append(
		s.limit.DaysOfWeek,
		time.Saturday,
	)
	return s
}

// Sundays limit the task to Sunday
func (s *Scheduler) Sundays() *Scheduler {
	s.limit.DaysOfWeek = append(
		s.limit.DaysOfWeek,
		time.Sunday,
	)
	return s
}

// Days limit the task to specific days
func (s *Scheduler) Days(d ...time.Weekday) *Scheduler {
	s.limit.DaysOfWeek = append(s.limit.DaysOfWeek, d...)
	return s
}

// Between limit the task to run between start and end time
func (s *Scheduler) Between(start, end string) *Scheduler {
	s.limit.StartTime = start
	s.limit.EndTime = end
	s.limit.IsBetween = true
	return s
}

// UnlessBetween limit the task to not run between start and end time
func (s *Scheduler) UnlessBetween(start, end string) *Scheduler {
	s.limit.StartTime = start
	s.limit.EndTime = end
	s.limit.IsBetween = false
	return s
}

// When limit the task based on a truth test
func (s *Scheduler) When(when WhenFunc) *Scheduler {
	s.limit.When = when
	return s
}
