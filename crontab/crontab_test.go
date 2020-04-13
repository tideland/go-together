// Tideland Go Together - CronTab - Unit Tests
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crontab_test

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/crontab"
)

//--------------------
// TESTS
//--------------------

// TestSubmitStatusRevoke tests a simple submitting, status retrieval,
// and revoking.
func TestSubmitStatusRevoke(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	at := time.Now().Add(250 * time.Millisecond)
	beenThereDoneThat := false

	// Test.
	err := crontab.SubmitAt("foo-1", at, func() error {
		beenThereDoneThat = true
		return nil
	})
	assert.NoError(err)
	err = crontab.Revoke("foo-1")
	assert.NoError(err)
	time.Sleep((500 * time.Millisecond))
	assert.False(beenThereDoneThat)

	err = crontab.SubmitAt("foo-2", at, func() error {
		return errors.New("ouch")
	})
	assert.NoError(err)
	time.Sleep(time.Second)
	err = crontab.Revoke("foo-2")
	assert.ErrorMatch(err, "ouch")
}

// TestList tests the listing of submitted jobs.
func TestList(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	at := time.Now().Add(time.Second)

	// Test.
	jobs, err := crontab.List()
	assert.NoError(err)
	assert.Empty(jobs)

	err = crontab.SubmitAt("fuddel-1", at, func() error {
		return nil
	})
	assert.NoError(err)
	jobs, err = crontab.List()
	assert.NoError(err)
	assert.Length(jobs, 1)

	err = crontab.SubmitAt("fuddel-2", at, func() error {
		return nil
	})
	assert.NoError(err)
	jobs, err = crontab.List()
	assert.NoError(err)
	assert.Length(jobs, 2)
	assert.Contains("fuddel-1", jobs)
	assert.Contains("fuddel-2", jobs)

	err = crontab.Revoke("fuddel-1")
	assert.NoError(err)
	jobs, err = crontab.List()
	assert.NoError(err)
	assert.Length(jobs, 1)
	assert.False(func() bool {
		for _, job := range jobs {
			if job == "fuddel-1" {
				return true
			}
		}
		return false
	}())
	assert.Contains("fuddel-2", jobs)

	err = crontab.Revoke("fuddel-2")
	assert.NoError(err)
	jobs, err = crontab.List()
	assert.NoError(err)
	assert.Empty(jobs)
}

// TestSubmitAt tests if a job is executed only once.
func TestSubmitAt(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	atOne := time.Now().Add(100 * time.Millisecond)
	oneDiffC := make(chan time.Duration, 1)
	oneDiffExtend := float64(10 * time.Millisecond)
	atTwo := atOne.Add(200 * time.Millisecond)
	syncC := make(chan struct{}, 10)

	// Test.
	err := crontab.SubmitAt("bar-1", atOne, func() error {
		oneDiffC <- time.Since(atOne)
		syncC <- struct{}{}
		return nil
	})
	assert.NoError(err)
	err = crontab.SubmitAt("bar-2", atTwo, func() error {
		syncC <- struct{}{}
		return nil
	})
	assert.NoError(err)

	count := 0
waiting:
	for {
		select {
		case <-syncC:
			count++
		case <-time.After(time.Second):
			break waiting
		}
	}
	assert.About(float64(<-oneDiffC), 0.0, oneDiffExtend)
	assert.Equal(count, 2)

	err = crontab.Revoke("bar-1")
	assert.NoError(err)
	err = crontab.Revoke("bar-2")
	assert.NoError(err)
}

// TestSubmitEvery tests if a job is executed every given interval.
func TestSubmitEvery(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	every := 200 * time.Millisecond
	syncC := make(chan struct{}, 1)

	// Test.
	err := crontab.SubmitEvery("baz-1", every, func() error {
		syncC <- struct{}{}
		return nil
	})
	assert.NoError(err)

	count := 0
	start := time.Now()
	for range syncC {
		count++
		if count > 10 {
			break
		}
	}
	duration := time.Since(start)

	assert.Range(duration, 2100*time.Millisecond, 2300*time.Millisecond)

	err = crontab.Revoke("baz-1")
	assert.NoError(err)
}

// TestSubmitAtEvery tests if a job is executed every given interval
// after a given time.
func TestSubmitAtEvery(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	at := time.Now().Add(500 * time.Millisecond)
	every := 100 * time.Millisecond
	syncC := make(chan struct{}, 1)

	// Test.
	err := crontab.SubmitAtEvery("babbel-1", at, every, func() error {
		syncC <- struct{}{}
		return nil
	})
	assert.NoError(err)

	count := 0
	start := time.Now()
	for range syncC {
		count++
		if count > 10 {
			break
		}
	}
	duration := time.Since(start)

	assert.Range(duration, 1500*time.Millisecond, 1700*time.Millisecond)

	err = crontab.Revoke("babbel-1")
	assert.NoError(err)
}

// TestSubmitAfterEvery tests if a job is executed every given interval
// after a given pause.
func TestSubmitAfterEvery(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	pause := 500 * time.Millisecond
	every := 100 * time.Millisecond
	syncC := make(chan struct{}, 1)

	// Test.
	err := crontab.SubmitAfterEvery("daddel-1", pause, every, func() error {
		syncC <- struct{}{}
		return nil
	})
	assert.NoError(err)

	count := 0
	start := time.Now()
	for range syncC {
		count++
		if count > 10 {
			break
		}
	}
	duration := time.Since(start)

	assert.Range(duration, 1500*time.Millisecond, 1700*time.Millisecond)

	err = crontab.Revoke("daddel-1")
	assert.NoError(err)
}

// TestIllegal tests double id submit and illegal id revoke.
func TestIllegal(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, asserts.FailStop)
	at := time.Now().Add(time.Second)
	job := func() error {
		return nil
	}

	// Test.
	err := crontab.SubmitAt("yadda-1", at, job)
	assert.NoError(err)
	err = crontab.SubmitAt("yadda-1", at, job)
	assert.ErrorMatch(err, `job ID 'yadda-1' already exists`)

	err = crontab.Revoke("yadda-2")
	assert.ErrorMatch(err, `job ID 'yadda-2' does not exist`)
}

// EOF
