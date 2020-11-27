# Tideland Go Together

[![GitHub release](https://img.shields.io/github/release/tideland/go-together.svg)](https://github.com/tideland/go-together)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-together/master/LICENSE)
[![Go Module](https://img.shields.io/github/go-mod/go-version/tideland/go-together)](https://github.com/tideland/go-together/blob/master/go.mod)
[![GoDoc](https://godoc.org/tideland.dev/go/together?status.svg)](https://pkg.go.dev/mod/tideland.dev/go/together?tab=packages)
[![Workflow](https://img.shields.io/github/workflow/status/tideland/go-together/build)](https://github.com/tideland/go-together/actions/)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/go-together)](https://goreportcard.com/report/tideland.dev/go/together)

## Description

**Tideland Go Together** focusses on goroutines and how to manage them more convenient and reliable.

* `actor` runs a backend goroutine processing anonymous functions for the serialization of changes, e.g. in a structure
* `cells` provides an event processing based on the idea of meshed cells with different behaviors
* `fuse` contains some ways of status and error control in concurrent applications
* `limiter` limits the number of parallel executing goroutines in its scope
* `loop` helps running a controlled endless `select` loop for goroutine backends
* `wait` provides a flexible and controlled waiting for conditions by polling

I hope you like it. ;)

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://tideland.dev)

