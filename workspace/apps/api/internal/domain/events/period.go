package events

import "strings"

// Period represents the match period when an event occurred.
type Period string

const (
	// PeriodFirstHalf represents the first half (0-45 minutes).
	PeriodFirstHalf Period = "first_half"
	// PeriodSecondHalf represents the second half (45-90 minutes).
	PeriodSecondHalf Period = "second_half"
	// PeriodExtraTimeFirst represents the first half of extra time (90-105 minutes).
	PeriodExtraTimeFirst Period = "extra_time_first"
	// PeriodExtraTimeSecond represents the second half of extra time (105-120 minutes).
	PeriodExtraTimeSecond Period = "extra_time_second"
	// PeriodPenalties represents the penalty shootout.
	PeriodPenalties Period = "penalties"
	// PeriodRegular is a fallback for regular time (will be auto-determined from minute).
	PeriodRegular Period = "regular"
)

// String returns the string representation of the period.
func (p Period) String() string {
	return string(p)
}

// IsValid checks if the period is valid.
func (p Period) IsValid() bool {
	switch p {
	case PeriodFirstHalf, PeriodSecondHalf, PeriodExtraTimeFirst,
		PeriodExtraTimeSecond, PeriodPenalties, PeriodRegular:
		return true
	default:
		return false
	}
}

// DeterminePeriod determines the period based on minute and extra_minute.
// This is useful when period is not explicitly provided.
func DeterminePeriod(minute int32, extraMinute *int32) Period {
	// If extra_minute is set, we're in extra time
	if extraMinute != nil && *extraMinute > 0 {
		if minute <= 105 {
			return PeriodExtraTimeFirst
		}
		return PeriodExtraTimeSecond
	}

	// Regular time
	if minute <= 45 {
		return PeriodFirstHalf
	}
	if minute <= 90 {
		return PeriodSecondHalf
	}

	// Beyond 90 minutes without extra_minute might be penalties or data error
	// Default to second half for safety
	return PeriodSecondHalf
}

// NormalizePeriod normalizes period strings from external providers.
// Handles various formats: "1H", "2H", "ET1", "ET2", "P", "first_half", etc.
func NormalizePeriod(periodStr string) Period {
	normalized := strings.ToLower(strings.TrimSpace(periodStr))

	switch normalized {
	case "1h", "1st_half", "first_half", "firsthalf", "first":
		return PeriodFirstHalf
	case "2h", "2nd_half", "second_half", "secondhalf", "second":
		return PeriodSecondHalf
	case "et1", "extra_time_1", "extra_time_first", "extratimefirst":
		return PeriodExtraTimeFirst
	case "et2", "extra_time_2", "extra_time_second", "extratimesecond":
		return PeriodExtraTimeSecond
	case "p", "pen", "penalties", "penalty_shootout", "shootout":
		return PeriodPenalties
	default:
		return PeriodRegular
	}
}
