package backoff

import (
	"testing"
	"time"
)

func TestForAttempt(t *testing.T) {
	jitterPct := 15.0

	b := Exponential(Config{
		Min:       50 * time.Millisecond,
		Max:       2 * time.Second,
		Factor:    4,
		JitterPct: jitterPct,
	})

	tests := []struct {
		name         string
		attempt      int
		wantNoJitter time.Duration
	}{
		{
			name:         "first attempt",
			attempt:      0,
			wantNoJitter: 50 * time.Millisecond,
		},
		{
			name:         "second attempt",
			attempt:      1,
			wantNoJitter: 200 * time.Millisecond,
		},
		{
			name:         "max overflow",
			attempt:      100,
			wantNoJitter: 2 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := b.Attempt(tt.attempt)

			dev := float64(got-tt.wantNoJitter) / float64(tt.wantNoJitter)
			if dev < 0 {
				dev *= -1
			}

			if dev > jitterPct/100 {
				t.Errorf("exponential.ForAttempt() = %v, want %v (Deviation %.6f > %.6f) seed %d", got, tt.wantNoJitter, dev, jitterPct/100, seed)
			}
		})
	}
}
