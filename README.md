# Tideland Go Together

[![GitHub release](https://img.shields.io/github/release/tideland/go-together.svg)](https://github.com/tideland/go-together)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-together/master/LICENSE)

## Description

**Tideland Go Together** focusses on goroutines and how to manage them more convenient and reliable.

* `actor` runs a backend goroutine processing anonymous functions for the serialization of changes, e.g. in a structure
* `cells` provides an event processing based on the idea of meshed cells with different behaviors
* `crontab` allows running functions at configured times and in chronological order
* `limiter` limits the number of parallel executing goroutines in its scope
* `loop` helps running a controlled endless `select` loop for goroutine backends
* `notifier` helps at the coordination of multiple goroutines
* `wait` provides a flexible and controlled waiting for conditions by polling

I hope you like it. ;)

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://tideland.dev)

## License

**Tideland Go Together** is distributed under the terms of the BSD 3-Clause license.
