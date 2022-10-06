# Schedule inspired laravel
A simple scheduler for golang, inspired with laravel's task scheduling, use it with `crontab`.

[![CI](https://github.com/iflamed/schedule/actions/workflows/ci.yml/badge.svg)](https://github.com/iflamed/schedule/actions/workflows/ci.yml) 
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/iflamed/schedule)](https://goreportcard.com/report/github.com/iflamed/schedule) 
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/iflamed/schedule/master/LICENSE) 
[![PkgGoDev](https://pkg.go.dev/badge/github.com/iflamed/schedule)](https://pkg.go.dev/github.com/iflamed/schedule)

### Run with crontab 
After you build your project, you can use it with crontab like below.
```shell
* * * * * /path/to/your/schedule >> /dev/null 2>&1
```

### Features
- You can integrate to your self project;
- Multiple frequency options;
- Day and time constraints;
- Custom logger;
- `context.Context` integration;
- Panic recover;
- Run task in go routine;
- Custom timezone;
- 100% code coverage;

### Installation
```shell
go get -u github.com/iflamed/schedule
```

## Getting Started
### Custom logger
**Logger interface**
```go
type Logger interface {
	Error(msg string, e any)
	Debugf(msg string, n int32)
	Debug(msg string)
}
```
**The default logger**
```go
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
```

### Create scheduler instance
```go
s := NewScheduler(context.Background(), time.UTC)
s.Daily().CallFunc(func(ctx context.Context) {
    log.Println("Task finished.")
    return
})
s.DailyAt("09:00").Call(NewDefaultTask(func(ctx context.Context) {
    log.Println("Task finished at 09:00")
}))
s.Start()
```
⚠️You should set frequency first, then call the `Call` and `CallFunc` method to run task.

⚠️`s.Start()` method must call at last, it will wait all task finished when process exit.

### Schedule Frequency Options
There are many more task schedule frequencies that you may assign to a task:

Method  | Description
------------- | -------------
`EveryMinute()`  |  Run the task every minute
`EveryTwoMinutes()`  |  Run the task every two minutes
`EveryThreeMinutes()`  |  Run the task every three minutes
`EveryFourMinutes()`  |  Run the task every four minutes
`EveryFiveMinutes()`  |  Run the task every five minutes
`EveryTenMinutes()`  |  Run the task every ten minutes
`EveryFifteenMinutes()`  |  Run the task every fifteen minutes
`EveryThirtyMinutes()`  |  Run the task every thirty minutes
`Hourly()`  |  Run the task every hour
`HourlyAt(17)`  |  Run the task every hour at 17 minutes past the hour
`EveryOddHour()`  |  Run the task every odd hour
`EveryTwoHours()`  |  Run the task every two hours
`EveryThreeHours()`  |  Run the task every three hours
`EveryFourHours()`  |  Run the task every four hours
`EverySixHours()`  |  Run the task every six hours
`Daily()`  |  Run the task every day at midnight
`DailyAt("13:00")`  |  Run the task every day at 13:00
`At("13:00")`  |  Run the task every day at 13:00, method alias of `dailyAt`
`TwiceDaily(1, 13)`  |  Run the task daily at 1:00 & 13:00
`TwiceDailyAt(1, 13, 15)`  |  Run the task daily at 1:15 & 13:15
`Weekly()`  |  Run the task every Sunday at 00:00
`WeeklyOn(1, "8:00")`  |  Run the task every week on Monday at 8:00
`Monthly()`  |  Run the task on the first day of every month at 00:00
`MonthlyOn(4, "15:00")`  |  Run the task every month on the 4th at 15:00
`TwiceMonthly(1, 16, "13:00")`  |  Run the task monthly on the 1st and 16th at 13:00
`LastDayOfMonth("15:00")` | Run the task on the last day of the month at 15:00
`Quarterly()` |  Run the task on the first day of every quarter at 00:00
`Yearly()`  |  Run the task on the first day of every year at 00:00
`YearlyOn(6, 1, "17:00")`  |  Run the task every year on June 1st at 17:00
`Timezone(time.UTC)` | Set the timezone for the task

### Schedule constraints
Method  | Description
------------- | -------------
`Weekdays()`  |  Limit the task to weekdays
`Weekends()`  |  Limit the task to weekends
`Sundays()`  |  Limit the task to Sunday
`Mondays()`  |  Limit the task to Monday
`Tuesdays()`  |  Limit the task to Tuesday
`Wednesdays()`  |  Limit the task to Wednesday
`Thursdays()`  |  Limit the task to Thursday
`Fridays()`  |  Limit the task to Friday
`Saturdays()`  |  Limit the task to Saturday
`Days(d ...time.Weekday)`  |  Limit the task to specific days
`Between(start, end string)`  |  Limit the task to run between start and end times
`UnlessBetween(start, end string)`  |  Limit the task to not run between start and end times
`When(when WhenFunc)`  |  Limit the task based on a truth test

### Schedule example
```go
package main

import (
	"context"
	"github.com/iflamed/schedule"
	"log"
	"time"
)

func main() {
	s := schedule.NewScheduler(context.Background(), time.UTC)
	s.Daily().CallFunc(func(ctx context.Context) {
		log.Println("Task finished.")
	})
	s.DailyAt("09:00").Call(schedule.NewDefaultTask(func(ctx context.Context) {
		log.Println("Task finished at 09:00")
	}))
	s.EveryMinute().Sundays().Call(schedule.NewDefaultTask(func(ctx context.Context) {
		log.Println("Task finished at 09:00")
	}))
	s.EveryMinute().Sundays().Between("12:00", "20:00").Call(schedule.NewDefaultTask(func(ctx context.Context) {
		log.Println("Task finished at 09:00")
	}))
	s.Start()
}

```
