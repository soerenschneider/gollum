package requeue

import (
	"math"
	"math/rand/v2"
	"time"
)

const baseBackoff = time.Minute * 2

// JitterPercentageAdditive adds random jitter to a given duration by a specified percentage.
// The jitter is additive, meaning the resulting duration is the original duration plus a random amount
// based on the given percentage of the duration.
//
// Parameters:
//   - duration: The base duration to apply jitter to.
//   - percentage: The percentage of jitter to add to the duration (e.g., 10 for 10%).
//
// Returns:
//   - A new duration with the added jitter.
//
// Example:
//   - If the duration is 100 seconds and the percentage is 10, the result will be a duration between
//     100 and 110 seconds.
func JitterPercentageAdditive(duration time.Duration, percentage float64) time.Duration {
	jitterPercent := rand.Float64() * percentage
	jitterAmount := duration.Seconds() * jitterPercent / 100
	return time.Duration(jitterAmount+duration.Seconds()) * time.Second
}

// JitterPercentageDistributed adds random jitter to a given duration by a specified percentage,
// but the jitter is distributed both positively and negatively (i.e., +/- jitter).
// This means the resulting duration could be shorter or longer than the original duration,
// depending on the randomly generated jitter.
//
// Parameters:
//   - duration: The base duration to apply jitter to.
//   - percentage: The maximum percentage of jitter to add or subtract (e.g., 10 for +/-10%).
//
// Returns:
//   - A new duration with the distributed jitter.
//
// Example:
//   - If the duration is 100 seconds and the percentage is 10, the result will be a duration between
//     90 and 110 seconds.
func JitterPercentageDistributed(duration time.Duration, percentage float64) time.Duration {
	jitterPercent := rand.Float64()*2*percentage - percentage
	jitterAmount := duration.Seconds() * jitterPercent / 100
	return time.Duration(jitterAmount+duration.Seconds()) * time.Second
}

// JitterFixAdditive adds a fixed random jitter to a given duration, with the amount of jitter
// determined by a random number of minutes between 0 and the specified maximum (inclusive).
// The jitter is additive, meaning the resulting duration is the original duration plus the random jitter.
//
// Parameters:
//   - duration: The base duration to apply jitter to.
//   - maxMinutes: The maximum number of minutes to add as jitter (must be greater than 0).
//
// Returns:
//   - A new duration with the added jitter, measured in minutes.
//
// Example:
//   - If the duration is 30 minutes and maxMinutes is 5, the result will be between 30 and 35 minutes.
func JitterFixAdditive(duration time.Duration, maxMinutes int) time.Duration {
	if maxMinutes <= 0 {
		maxMinutes = 5
	}
	jitterAmount := duration.Minutes() + float64(rand.IntN(maxMinutes))
	return time.Duration(jitterAmount) * time.Minute
}

// ExponentialBackoff calculates the backoff duration for retries using an exponential backoff
// strategy. The backoff time increases exponentially with each retry, up to a minimum threshold.
//
// Parameters:
//   - retries: The number of retries that have been attempted (starting from 0).
//   - minWaitMinutes: The minimum wait time in minutes before retrying, which is the lower bound
//     of the backoff calculation.
//
// Returns:
//   - The calculated backoff duration based on the number of retries and the minimum wait time.
//
// Example:
//   - If retries is 2 and minWaitMinutes is 5, the backoff time will be the greater of
//     5 minutes or 2^2 * baseBackoff (e.g., 4 * baseBackoff).
func ExponentialBackoff(retries int, minWaitMinutes int) time.Duration {
	minWait := time.Duration(minWaitMinutes) * time.Minute
	backoffTime := baseBackoff * time.Duration(math.Pow(2, float64(retries)))
	if backoffTime < minWait {
		return time.Duration(minWaitMinutes) * time.Minute
	}

	return backoffTime
}
