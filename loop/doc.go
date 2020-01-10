// Tideland Go Together - Loop
//
// Copyright (C) 2017-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package loop supports the developer implementing the typical Go
// idiom for concurrent applications running in a loop in the
// background and doing a select on one or more channels. Stopping
// those loops or getting aware of internal errors requires extra
// efforts. The loop package helps to control this kind of goroutines.
//
//     type Printer struct {
//	       printC chan string
//         loop   loop.Loop
//     }
//
//     func NewPrinter() (*Printer, error) {
//         p := &Printer{
//             printC: make(chan string),
//         }
//         l, err := loop.Go(p.worker)
//         if err != nil {
//             return nil, err
//         }
//         p.loop = l
//         return p, nil
//     }
//
//     func (p *printer) worker(c *notifier.Closer) error {
//         for {
//             select {
//             case <-c.Done():
//                 return nil
//             case str := <-printC:
//                 println(str)
//         }
//     }
//
// The worker here now can be stopped with p.loop.Stop(err) returning
// a possible internal error. Also recovering of internal errors or
// panics by starting the loop with a recoverer function is possible.
// See the code examples.
package loop // import "tideland.dev/go/together/loop"

// EOF
