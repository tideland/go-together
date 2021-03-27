// Tideland Go Together - Wait - Unit Tests
//
// Copyright (C) 2019-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package wait_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/together/wait"
)

//--------------------
// TESTS
//--------------------

// TestPollWithChangingInterval tests the polling of conditions with
// changing interval durations.
func TestPollWithChangingInterval(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	makeChanger := func(interval time.Duration) wait.TickChanger {
		return func(in time.Duration) (out time.Duration, ok bool) {
			if in == 0 {
				return interval, true
			}
			out = in * 2
			if out > time.Second {
				return 0, false
			}
			return out, true
		}
	}

	// Tests.
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeGenericIntervalTicker(makeChanger(10*time.Millisecond)),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeGenericIntervalTicker(makeChanger(10*time.Millisecond)),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Equal(count, 7, "exceeded with a count")

	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeGenericIntervalTicker(makeChanger(10*time.Millisecond)),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*cancelled.*")
	assert.Range(count, 4, 6, "test is race, depending on scheduling")
}

// TestPollWithInterval tests the polling of conditions in intervals.
func TestPollWithInterval(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Tests.
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeIntervalTicker(20*time.Millisecond),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	count = 0
	err = wait.WithInterval(context.Background(), 20*time.Millisecond, func() (bool, error) {
		count++
		if count == 5 {
			return true, nil
		}
		return false, nil
	})
	assert.NoError(err)
	assert.Equal(count, 5)

	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeIntervalTicker(20*time.Millisecond),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*cancelled.*")
	assert.Range(count, 4, 6, "test is race, depending on scheduling")
}

// TestPollWithMaxIntervals tests the polling of conditions in a maximum
// number of intervals.
func TestPollWithMaxInterval(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Tests.
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeMaxIntervalsTicker(20*time.Millisecond, 10),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	count = 0
	err = wait.WithMaxIntervals(context.Background(), 20*time.Millisecond, 10, func() (bool, error) {
		count++
		if count == 5 {
			return true, nil
		}
		return false, nil
	})
	assert.NoError(err)
	assert.Equal(count, 5)

	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeMaxIntervalsTicker(20*time.Millisecond, 10),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Equal(count, 10, "exceeded with a count")

	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeMaxIntervalsTicker(20*time.Millisecond, -1),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Equal(count, 0)

	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeMaxIntervalsTicker(20*time.Millisecond, 10),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*cancelled.*")
	assert.Range(count, 4, 6, "test is race, depending on scheduling")
}

// TestPollWithDeadline tests the polling of conditions with deadlines.
func TestPollWithDeadline(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Tests.
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeDeadlinedIntervalTicker(20*time.Millisecond, time.Now().Add(210*time.Millisecond)),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	count = 0
	err = wait.WithDeadline(context.Background(), 20*time.Millisecond, time.Now().Add(210*time.Millisecond), func() (bool, error) {
		count++
		if count == 5 {
			return true, nil
		}
		return false, nil
	})
	assert.NoError(err)
	assert.Equal(count, 5)

	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeDeadlinedIntervalTicker(100*time.Millisecond, time.Now().Add(1020*time.Millisecond)),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Range(count, 9, 11, "exceeded with a count")

	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeDeadlinedIntervalTicker(20*time.Millisecond, time.Now().Add(-time.Second)),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Equal(count, 0)

	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeDeadlinedIntervalTicker(20*time.Millisecond, time.Now().Add(time.Second)),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*cancelled.*")
	assert.Range(count, 4, 6, "test is race, depending on scheduling")
}

// TestPollWithTimeout tests the polling of conditions with timeouts.
func TestPollWithTimeout(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Tests.
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeExpiringIntervalTicker(20*time.Millisecond, 210*time.Millisecond),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	count = 0
	err = wait.WithTimeout(context.Background(), 20*time.Millisecond, 210*time.Millisecond, func() (bool, error) {
		count++
		if count == 5 {
			return true, nil
		}
		return false, nil
	})
	assert.NoError(err)
	assert.Equal(count, 5)

	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeExpiringIntervalTicker(20*time.Millisecond, 210*time.Millisecond),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Range(count, 9, 11)

	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeExpiringIntervalTicker(20*time.Millisecond, -10*time.Millisecond),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Equal(count, 0)

	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeExpiringIntervalTicker(20*time.Millisecond, 500*time.Millisecond),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*cancelled.*")
	assert.Range(count, 4, 6, "test is race, depending on scheduling")
}

// TestPollWithJitter tests the polling of conditions in a maximum
// number of intervals.
func TestPollWithJitter(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Tests.
	timestamps := []time.Time{}
	err := wait.Poll(
		context.Background(),
		wait.MakeJitteringTicker(50*time.Millisecond, 1.0, 1250*time.Millisecond),
		func() (bool, error) {
			timestamps = append(timestamps, time.Now())
			if len(timestamps) == 10 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Length(timestamps, 10)
	for i := 1; i < 10; i++ {
		diff := timestamps[i].Sub(timestamps[i-1])
		// 10% upper tolerance.
		assert.Range(diff, 50*time.Millisecond, 110*time.Millisecond)
	}

	timestamps = []time.Time{}
	err = wait.WithJitter(context.Background(), 50*time.Millisecond, 1.0, 1250*time.Millisecond, func() (bool, error) {
		timestamps = append(timestamps, time.Now())
		if len(timestamps) == 10 {
			return true, nil
		}
		return false, nil
	})
	assert.NoError(err)
	assert.Length(timestamps, 10)
	for i := 1; i < 10; i++ {
		diff := timestamps[i].Sub(timestamps[i-1])
		// 10% upper tolerance.
		assert.Range(diff, 50*time.Millisecond, 110*time.Millisecond)
	}

	timestamps = []time.Time{}
	err = wait.Poll(
		context.Background(),
		wait.MakeJitteringTicker(50*time.Millisecond, 1.0, 1250*time.Millisecond),
		func() (bool, error) {
			timestamps = append(timestamps, time.Now())
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Range(len(timestamps), 10, 25)

	timestamps = []time.Time{}
	err = wait.Poll(
		context.Background(),
		wait.MakeJitteringTicker(50*time.Millisecond, 1.0, -10*time.Millisecond),
		func() (bool, error) {
			timestamps = append(timestamps, time.Now())
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Empty(timestamps)

	timestamps = []time.Time{}
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeJitteringTicker(50*time.Millisecond, 1.0, 1250*time.Millisecond),
		func() (bool, error) {
			timestamps = append(timestamps, time.Now())
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*cancelled.*")
	assert.Range(len(timestamps), 3, 7, "test is race, depending on scheduling")
}

// TestPoll tests the polling of conditions with a user-defined ticker.
func TestPoll(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	ticker := func(ctx context.Context) <-chan struct{} {
		// Ticker runs 1000 times.
		tickc := make(chan struct{})
		go func() {
			count := 0
			defer close(tickc)
			for {
				select {
				case tickc <- struct{}{}:
					count++
					if count == 1000 {
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
		return tickc
	}

	// Tests.
	count := 0
	err := wait.Poll(
		context.Background(),
		ticker,
		func() (bool, error) {
			count++
			if count == 500 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 500)

	count = 0
	err = wait.Poll(
		context.Background(),
		ticker,
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*exceeded.*")
	assert.Equal(count, 1000, "exceeded with a count")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		ticker,
		func() (bool, error) {
			time.Sleep(2 * time.Millisecond)
			return false, nil
		},
	)
	assert.ErrorMatch(err, ".*cancelled.*")
}

// TestPanic tests the handling of panics during condition checks.
func TestPanic(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)

	// Test.
	count := 0
	err := wait.WithInterval(context.Background(), 10*time.Millisecond, func() (bool, error) {
		count++
		if count == 5 {
			panic("ouch at five o'clock")
		}
		return false, nil
	})
	assert.ErrorMatch(err, ".*panic.*")
	assert.Equal(count, 5)
}

// EOF
