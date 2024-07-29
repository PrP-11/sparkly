package pkg

import "time"

var AnalyticsTimeFrames map[string]time.Duration = map[string]time.Duration{
	"last_minute": time.Minute,
	"last_hour":   time.Hour,
	"last_day":    24 * time.Hour,
}
