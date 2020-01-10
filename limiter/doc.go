// Tideland Go Together - Limiter
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package limiter allows to restrict the number of parallel executed
// functions, e.g. to avoid an overload where processing requests.
//
//     l := limiter.New(10)
//     job := func() { ... } error
//
//     for i := 0; i < 100; i++ {
//         go func() {
//             err := l.Do(ctx, job)
//         }()
//     }
//
// Here even with a large number of goroutines the execution of function job
// is restricted to 10 at the same time.
package limiter // import "tideland.dev/go/together/limiter"

// EOF
