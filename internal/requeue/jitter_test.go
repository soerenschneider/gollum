package requeue

import (
	"testing"
	"time"
)

const (
	defaultRequeueInterval          = time.Hour
	defaultJitterPercentage float64 = 20
)

func TestJitterPercentageDistributed(t *testing.T) {
	for i := 0; i < 2000; i++ {
		got := JitterPercentageDistributed(defaultRequeueInterval, defaultJitterPercentage)
		if got < time.Duration(40)*time.Minute || got > time.Duration(80)*time.Minute {
			t.Errorf("expected [40, 80], got %v", got)
		}
	}
}
