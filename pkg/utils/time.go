package utils

import (
	"fmt"
	"time"
)

// RFC3339Time is an alias of time.RFC3339Time.
type RFC3339Time time.Time

// MarshalJSON marshals the response.
func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	ts := time.Time(t)
	return []byte(fmt.Sprintf("%q", ts.Format(time.RFC3339))), nil
}

func (t RFC3339Time) String() string {
	ts := time.Time(t)
	return ts.Format(time.RFC3339)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time must be a quoted string in the RFC 3339 format.
func (t *RFC3339Time) UnmarshalJSON(data []byte) error {
	var tt time.Time
	if err := tt.UnmarshalJSON(data); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	*t = RFC3339Time(tt)

	return nil
}

// FirstDayOfWeek calculates the first day of the week (Monday) for a given date.
func FirstDayOfWeek(t time.Time) time.Time {
	// Calculate the day difference from Monday (0 = Monday, ..., 6 = Sunday)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Adjust Sunday to 7 for easier math
	}

	// Subtract the number of days to get back to Monday
	t = t.AddDate(0, 0, -weekday+1)

	// rm time from date
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return t
}

const (
	// WeekDaysCount is count of week days.
	WeekDaysCount = 7
)
