// Package backoff provides exponential backoff duration generator.
package backoff

import (
	"math/rand"
	"sync"
	"time"
)

var (
	seed              = time.Now().UnixNano()
	randGenerator     = rand.New(rand.NewSource(seed))
	randGeneratorLock sync.Mutex
)

var defaultConf = Config{
	Min:       100 * time.Millisecond,
	Max:       3 * time.Second,
	JitterPct: 20,
	Factor:    2,
}

// Config - backoff configuration options
type Config struct {
	// Min - backoff's initial duration.
	// Default - 100ms
	Min time.Duration
	// Max - maximal duration of backoff. When calculated backoff is higher than Max, Max is returned.
	// Default - 3s
	Max time.Duration
	// Factor - backoff multiplication factor between the attempts. Should be > 1
	// Default - 2
	Factor float64
	// JitterPct - backoff deviation percent.
	// Set to negative number, if jitter should be disabled
	// Default - 20%
	JitterPct float64
}

// Backoff - generator of backoff intervals
type Backoff interface {
	Attempt(num int) time.Duration
}

type exponential struct {
	min    float64
	max    float64
	factor float64
	jitter float64
}

// Exponential - creates instance of exponential backoff generator. Is goroutine safe
func Exponential(conf Config) Backoff {

	var e exponential

	if conf.Min > 0 {
		e.min = float64(conf.Min)
	} else {
		e.min = float64(defaultConf.Min)
	}

	if conf.Max > 0 {
		e.max = float64(conf.Max)
	} else {
		e.max = float64(defaultConf.Max)
	}

	if conf.Factor > 1 {
		e.factor = conf.Factor
	} else {
		e.factor = defaultConf.Factor
	}

	if conf.JitterPct > 0 {
		e.jitter = conf.JitterPct / 100
	} else if conf.JitterPct == 0 {
		e.jitter = defaultConf.JitterPct / 100
	}

	return &e
}

func (e *exponential) jitterDeviation() float64 {
	randGeneratorLock.Lock()
	rf := randGenerator.Float64()
	randGeneratorLock.Unlock()

	return 1 + e.jitter*(2*rf-1)
}

// Attempt - returns backoff duration for Nth retry attempt
func (e *exponential) Attempt(attempt int) time.Duration {

	backoff := e.min
	for ; attempt > 0 && backoff <= e.max; attempt-- {
		backoff *= e.factor
	}
	if backoff > e.max {
		backoff = e.max
	}
	if e.jitter > 0 {
		backoff *= e.jitterDeviation()
		if backoff < 0 {
			return 0
		}
	}
	return time.Duration(backoff)
}
