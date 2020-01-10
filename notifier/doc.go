// Tideland Go Together - Notifier
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package notifier helps at the coordination of multiple goroutines. First little
// helper is the Closer for the aggregation of typical closer channels <-chan struct{},
// i.e. by Context.Done(), into one Closer.Done() <-chan struct{}. This way for-select-loops
// don't need to test each channel individually.
//
//     ca := make(chan struct{})
//     cb := make(chan struct{})
//     closer := notifier.NewCloser(ca, cb, ctx.Done())
//
//     go func() {
//          for {
//              select {
//              case <-closer.Done():
//                  return
//              case foo := <-someOtherChannel:
//                  ...
//              }
//          }
//     }()
//
//     close(cb)  // Or any other of the passed channels.
//
// Second and third type are Statebox and Bundle. The Statebox is intended to be
// in one of the states Unknown, Starting, Ready, Working, Stopping, and Stopped
// and is notifying an interested and waiting goroutine.
//
//     statebox := notifier.NewStatebox()
//
//     go func() {
//         statebox.Notify(notifier.Starting)
//         // Init goroutine.
//         ...
//         statebox.Notify(notifier.Ready)
//         ...
//         // After done work and other states.
//         statebox.Notify(notifier.Stopped)
//     }()
//
//     if err := statebox.Wait(notifier.Working, timeout); err != nil {
//         panic("timeout")
//     }
//     log.Printf("system is working")
//
//     // Do something too.
//     ...
//     if err := statebox.Wait(notifier.Stopped, timeout); err != nil {
//         panic("timeout")
//     }
//     log.Printf("system has stopped")
//
// The Bundle simply helps to notify multiple Stateboxes en bloc.
package notifier // import "tideland.dev/go/together/notifier"

// EOF
