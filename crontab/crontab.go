// Tideland Go Together - CronTab
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crontab // import "tideland.dev/go/together/crontab"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"tideland.dev/go/together/actor"
	"tideland.dev/go/together/fuse"
	"tideland.dev/go/together/loop"
)

//--------------------
// GLOBAL
//--------------------

var (
	mu sync.Mutex
	ct *crontab
)

//--------------------
// CRONTAB
//--------------------

// crontab implements the tanle for all cronjobs.
type crontab struct {
	act  *actor.Actor
	jobs map[string]*cronjob
}

// goCrontab instantiates and starts the crontab if it is not already running.
func goGrontab() {
	mu.Lock()
	defer mu.Unlock()
	if ct != nil {
		return
	}
	act, err := actor.Go()
	if err != nil {
		panic("starting crontab: " + err.Error())
	}
	ct = &crontab{
		act:  act,
		jobs: make(map[string]*cronjob),
	}
}

//--------------------
// CRONJOB
//--------------------

// cronjob is responsible to run one job.
type cronjob struct {
	id       string
	start    *time.Time
	interval *time.Duration
	job      func() error
	loop     *loop.Loop
	rs       fuse.Reasons
}

// newCronjob creates a new cronjob and starts its goroutine.
func newCronjob(id string, s *time.Time, i *time.Duration, j func() error) *cronjob {
	cj := &cronjob{
		id:       id,
		start:    s,
		interval: i,
		job:      j,
	}
	l, err := loop.Go(
		cj.worker,
		loop.WithRecoverer(cj.recoverer),
	)
	if err != nil {
		panic("start cronjob: " + err.Error())
	}
	cj.loop = l
	return cj
}

// stop ends the cronjob goroutine.
func (cj *cronjob) stop() error {
	return cj.loop.Stop()
}

// worker runs the cronjob.
func (cj *cronjob) worker(lt loop.Terminator) error {
	// Init.
	var interval time.Duration
	if cj.start != nil {
		interval = time.Until(*cj.start)
	} else {
		interval = *cj.interval
	}
	// Loop.
	for {
		select {
		case <-lt.Done():
			return nil
		case <-time.After(interval):
			if err := cj.job(); err != nil {
				return err
			}
			if cj.interval == nil {
				// Only once.
				return nil
			}
			// In intervals.
			interval = *cj.interval
		}
	}
}

// recoverer allows the cronjob to survive panics.
func (cj *cronjob) recoverer(reason interface{}) error {
	cj.rs.Append(reason)
	if cj.rs.Frequency(5, 10*time.Millisecond) {
		return fmt.Errorf("too high error frequency: %v", cj.rs)
	}
	if cj.rs.Len() >= 10 {
		return fmt.Errorf("too many errors: %v", cj.rs)
	}
	return nil
}

//--------------------
// API
//--------------------

// SubmitAt adds a function running only once at a given time.
func SubmitAt(id string, at time.Time, j func() error) error {
	goGrontab()
	var err error
	if aerr := ct.act.DoSync(func() {
		if ct.jobs[id] != nil {
			err = fmt.Errorf("job ID '%s' already exists", id)
			return
		}
		ct.jobs[id] = newCronjob(id, &at, nil, j)
	}); aerr != nil {
		return aerr
	}
	return err
}

// SubmitEvery adds a function running every interval.
func SubmitEvery(id string, every time.Duration, j func() error) error {
	goGrontab()
	var err error
	if aerr := ct.act.DoSync(func() {
		if ct.jobs[id] != nil {
			err = fmt.Errorf("job ID '%s' already exists", id)
			return
		}
		ct.jobs[id] = newCronjob(id, nil, &every, j)
	}); aerr != nil {
		return aerr
	}
	return err
}

// SubmitAtEvery adds a function running every interval starting at a given time.
func SubmitAtEvery(id string, at time.Time, every time.Duration, j func() error) error {
	goGrontab()
	var err error
	if aerr := ct.act.DoSync(func() {
		if ct.jobs[id] != nil {
			err = fmt.Errorf("job ID '%s' already exists", id)
			return
		}
		ct.jobs[id] = newCronjob(id, &at, &every, j)
	}); aerr != nil {
		return aerr
	}
	return err
}

// SubmitAfterEvery adds a function running every interval after a given pause.
func SubmitAfterEvery(id string, pause, every time.Duration, j func() error) error {
	return SubmitAtEvery(id, time.Now().Add(pause), every, j)
}

// List returns all currently submitted IDs.
func List() ([]string, error) {
	goGrontab()
	var ids []string
	var err error
	if aerr := ct.act.DoSync(func() {
		for id := range ct.jobs {
			ids = append(ids, id)
		}
	}); aerr != nil {
		return ids, aerr
	}
	sort.Strings(ids)
	return ids, err
}

// Revoke stops a cronjob and removes it from the table.
func Revoke(id string) error {
	goGrontab()
	var err error
	if aerr := ct.act.DoSync(func() {
		if job, ok := ct.jobs[id]; ok {
			delete(ct.jobs, id)
			err = job.stop()
			return
		}
		err = fmt.Errorf("job ID '%s' does not exist", id)
	}); aerr != nil {
		return aerr
	}
	return err
}

// EOF
