package clock

import (
	"time"
)

type Service interface {
	Now() time.Time
	LocalTime(timeZoneName string) time.Time
}

type clock struct{}

func NewClock() Service {
	return &clock{}
}

func (clock *clock) Now() time.Time {
	return time.Now()
}

func (clock *clock) LocalTime(timeZoneName string) time.Time {
	location, err := time.LoadLocation(timeZoneName)
	if err != nil {
		return clock.Now()
	}
	return clock.Now().In(location)
}
