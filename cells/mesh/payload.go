// Tideland Go Together - Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/together/cells/mesh"

//--------------------
// TOPICS
//--------------------

// Standard topics.
const (
	TopicNil        = "(nil)"
	TopicTerminated = "terminated"
	TopicError      = "error"
)

//--------------------
// PAYLOADS
//--------------------

// PayloadTermination describes the normal termination of a cell.
type PayloadTermination struct {
	CellName string `json:"cellName"`
}

// PayloadCellError describes the abnormal termination of a cell.
type PayloadCellError struct {
	CellName string `json:"cellName"`
	Error    string `json:"error"`
}

// EOF
