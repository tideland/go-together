// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2020 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"sort"

	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// EVALUATOR BEHAVIOR
//--------------------

// Evaluator is a function returning a rating for each received event.
type Evaluator func(evt *event.Event) (float64, error)

// evaluation contains the aggregated result of all evaluations.
type evaluation struct {
	count     int
	minRating float64
	maxRating float64
	avgRating float64
	medRating float64
}

// evaluatorBehavior implements the evaluator behavior.
type evaluatorBehavior struct {
	id            string
	emitter       mesh.Emitter
	evaluate      Evaluator
	maxRatings    int
	ratings       []float64
	sortedRatings []float64
	evaluation    evaluation
}

// NewEvaluatorBehavior creates a behavior evaluating received events based
// on the passed function. This function returns a rating. Their minimum,
// maximum, average, median, and number of events are emitted. The number
// of ratings for the median calculation is unlimited. So think about
// choosing NewMovingEvaluatorBehavior() to create the behavior with a
// limit and reduce memory usage.
//
// A "reset" topic sets all values to zero again.
func NewEvaluatorBehavior(id string, evaluator Evaluator) mesh.Behavior {
	return NewMovingEvaluatorBehavior(id, evaluator, 0)
}

// NewMovingEvaluatorBehavior creates the evaluator behavior with a
// moving rating window for calculation.
func NewMovingEvaluatorBehavior(id string, evaluator Evaluator, limit int) mesh.Behavior {
	return &evaluatorBehavior{
		id:         id,
		evaluate:   evaluator,
		maxRatings: limit,
		evaluation: evaluation{},
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *evaluatorBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *evaluatorBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *evaluatorBehavior) Terminate() error {
	b.ratings = nil
	b.sortedRatings = nil
	b.evaluation = evaluation{}
	return nil
}

// Process evaluates the event.
func (b *evaluatorBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case event.TopicReset:
		b.ratings = nil
		b.sortedRatings = nil
		b.evaluation = evaluation{}
	default:
		// Evaluate event and collect rating.
		rating, err := b.evaluate(evt)
		if err != nil {
			return err
		}
		b.ratings = append(b.ratings, rating)
		if b.maxRatings > 0 && len(b.ratings) > b.maxRatings {
			b.ratings = b.ratings[1:]
		}
		if len(b.sortedRatings) < len(b.ratings) {
			// Let it grow up to the needed size.
			b.sortedRatings = append(b.sortedRatings, 0.0)
		}
		// Evaluate ratings.
		b.evaluateRatings()
		b.emitter.Broadcast(event.New(
			TopicEvaluation,
			"count", b.evaluation.count,
			"min-rating", b.evaluation.minRating,
			"max-rating", b.evaluation.maxRating,
			"avg-rating", b.evaluation.avgRating,
			"med-rating", b.evaluation.medRating,
		))
	}
	return nil
}

// Recover from an error.
func (b *evaluatorBehavior) Recover(err interface{}) error {
	b.ratings = nil
	b.sortedRatings = nil
	b.evaluation = evaluation{}
	return nil
}

// evaluateRatings evaluates the collected ratings.
func (b *evaluatorBehavior) evaluateRatings() {
	copy(b.sortedRatings, b.ratings)
	sort.Float64s(b.sortedRatings)
	// Count.
	b.evaluation.count = len(b.sortedRatings)
	// Average.
	totalRating := 0.0
	for _, rating := range b.sortedRatings {
		totalRating += rating
	}
	b.evaluation.avgRating = totalRating / float64(b.evaluation.count)
	// Median.
	if b.evaluation.count%2 == 0 {
		// Even, have to calculate.
		middle := b.evaluation.count / 2
		b.evaluation.medRating = (b.sortedRatings[middle-1] + b.sortedRatings[middle]) / 2
	} else {
		// Odd, can take the middle.
		b.evaluation.medRating = b.sortedRatings[b.evaluation.count/2]
	}
	// Minimum and maximum.
	b.evaluation.minRating = b.sortedRatings[0]
	b.evaluation.maxRating = b.sortedRatings[len(b.sortedRatings)-1]
}

// EOF
