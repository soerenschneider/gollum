package requeue

import (
	"errors"
	"time"
)

type OnHoursRequeue struct {
	fromHour int
	toHour   int
}

// NewOnHoursRequeue creates a new OnHoursRequeue instance with the given time range
// specified by 'from' and 'to' hours. The function validates that the 'from' and 'to'
// values are within the valid range [0, 23] and ensures that 'from' is less than 'to'.
//
// Parameters:
//   - from: The starting hour of the on-hours period (0-23).
//   - to: The ending hour of the on-hours period (0-23).
//
// Returns:
//   - A pointer to an OnHoursRequeue struct initialized with the provided hours.
//   - An error if the provided hours are out of the valid range or if 'from' is not less than 'to'.
//
// Errors:
//   - "only values [0, 23] are allowed" if 'from' or 'to' are outside the valid range.
//   - "from must be < to" if 'from' is greater than or equal to 'to'.
func NewOnHoursRequeue(from, to int) (*OnHoursRequeue, error) {
	if from < 0 || to > 23 {
		return nil, errors.New("only values [0, 23] are allowed")
	}
	if from >= to {
		return nil, errors.New("from must be < to")
	}

	return &OnHoursRequeue{
		fromHour: from,
		toHour:   to,
	}, nil
}

func (w *OnHoursRequeue) Requeue(duration time.Duration) time.Duration {
	currentTime := time.Now().UTC().Add(duration)
	if currentTime.Hour() < w.fromHour {
		return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), w.fromHour, 0, 0, 0, time.UTC).Sub(currentTime)
	}

	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), w.fromHour, 0, 0, 0, time.UTC).AddDate(0, 0, 1).Sub(currentTime)
}
