// Tideland Go Together - CronTab
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package crontab provides the ability to run functions once at a given time,
// every given interval, every given interval after a given time, or every
// given interval after a given pause. The jobs are added by the different
// submit commands.
//
//     at := time.Date(2018, time.October, 31, 12, 0, 0, 0, time.UTC)
//     every := 24 * time.Hour
//
//     crontab.SubmitAt("id-1", at, func() error {
//         log.Printf("I'm executed once at %v.", time.Now())
//         return nil
//     })
//
//     crontab.SubmitEvery("id-2", every, func() error {
//         log.Printf("I'm executed every 24h. It is %v.", time.Now())
//         return nil
//     })
//
//     crontab.SubmitAtEvery("id-3", at, every, func() error {
//         log.Printf("I'm executed at %v first, every 24h. It is %v.", at, time.Now())
//         return nil
//     })
//
//     crontab.SubmitAfterEvery("id-4", time.Hour, every, func() error {
//         log.Printf("I'm executed first after an hour, then every %v.", every)
//         return nil
//     })
//
// Jobs can be deleted with crontab.Revoke(anyID) again.
package crontab // import "tideland.dev/go/together/crontab"

// EOF
