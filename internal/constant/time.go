package constant

import "time"

const (
	ThirtyDays   = 30 * 24 * time.Hour
	OneHour      = 1 * time.Hour
	TwoHours     = 2 * time.Hour
	OneDay       = 1 * 24 * time.Hour
	ThirtySecond = 30 * time.Second
	NinetySecond = 90 * time.Second

	TwoMinutes      = 2 * time.Minute
	ThreeMinutes    = 3 * time.Minute
	FiveMinutes     = 5 * time.Minute
	EightMinutes    = 8 * time.Minute
	TenMinutes      = 10 * time.Minute
	ThirteenMinutes = 13 * time.Minute
	FinalIncrement  = 5 * time.Minute

	MaxRetrySyncStatus = 5
)

var DelayBasedOnCounter = map[int]time.Duration{
	1: TwoMinutes,
	2: ThreeMinutes,
	3: FiveMinutes,
	4: EightMinutes,
	5: ThirteenMinutes,
}
