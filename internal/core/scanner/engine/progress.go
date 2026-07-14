package engine

import (
	"math"
	"time"
)

// Progress represents a thread-safe snapshot of the current execution status of the engine.
type Progress struct {
	Total      uint64        // Total number of tasks/IPs to process
	Processed  uint64        // Number of tasks already processed
	Succeed    uint64        // Number of successful tasks
	Percent    float64       // Completion percentage (0.0 to 100.0)
	Elapsed    time.Duration // Active elapsed time (excluding pause durations)
	RatePerSec float64       // Processing throughput rate (items/second)
	ETA        time.Duration // Estimated time remaining until completion
}

// maxETASeconds caps ETA to ~292 years, the largest value that fits in time.Duration (int64 nanoseconds).
const maxETASeconds = float64(math.MaxInt64) / float64(time.Second)

// reportProgress calculates current progress statistics and invokes the provided callback.
func reportProgress(
	start time.Time,
	paused time.Duration,
	total uint64,
	processed uint64,
	succeed uint64,
	cb func(p Progress),
) {
	if cb == nil {
		return
	}

	now := time.Now()
	elapsed := max(now.Sub(start)-paused, 0)

	var rate float64
	if elapsed > 0 {
		rate = float64(processed) / elapsed.Seconds()
	}

	var percent float64
	if total > 0 {
		percent = (float64(processed) / float64(total)) * 100.0
		if percent > 100.0 {
			percent = 100.0
		}
	}

	var eta time.Duration
	if rate > 0 && processed < total {
		etaSeconds := float64(total-processed) / rate
		if etaSeconds < maxETASeconds {
			eta = time.Duration(etaSeconds * float64(time.Second))
		} else {
			eta = time.Duration(math.MaxInt64)
		}
	}

	p := Progress{
		Total:      total,
		Processed:  processed,
		Succeed:    succeed,
		Percent:    percent,
		Elapsed:    elapsed,
		RatePerSec: rate,
		ETA:        eta,
	}

	cb(p)
}
