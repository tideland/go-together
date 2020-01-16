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
//	       prints chan string
//         loop   loop.Loop
//     }
//
//     func NewPrinter(ctx context.Context) (*Printer, error) {
//         p := &Printer{
//             ctx:    ctx,
//             prints: make(chan string),
//         }
//         l, err := loop.Go(
//             p.worker,
//             loop.WithContext(ctx),
//             loop.WithFinalizer(func(err error) error {
//                 ...
//             })
//         if err != nil {
//             return nil, err
//         }
//         p.loop = l
//         return p, nil
//     }
//
//     func (p *printer) worker(lt loop.Terminator) error {
//         for {
//             select {
//             case <-lt.Done():
//                 return nil
//             case str := <-p.prints:
//                 println(str)
//         }
//     }
//
// The worker here now can be stopped with p.loop.Stop() returning
// a possible internal error. Also recovering of internal errors or
// panics by starting the loop with a recoverer function is possible.
// See the code examples.
package loop // import "tideland.dev/go/together/loop"

// EOF
