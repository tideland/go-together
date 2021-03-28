// Tideland Go Together - Cells - Behaviors
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// TOPICS
//--------------------

// Standard topics.
const (
	TopicAggregate     = "aggregate!"
	TopicAggregated    = "aggregated"
	TopicCriterionDone = "criterion-done"
	TopicProcess       = "process!"
	TopicReset         = "reset!"
	TopicResetted      = "resetted"
)

// CriterionMatch allows a combo criterion func to signal its
// analysis rersult.
type CriterionMatch string

// Criterion matches.
const (
	CriterionDone      CriterionMatch = "criterion-done"
	CriterionKeep                     = "criterion-keep"
	CriterionDropFirst                = "criterion-drop-first"
	CriterionDropLast                 = "criterion-drop-last"
)

// EOF
