// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package tree

import (
	"github.com/cockroachdb/cockroach/pkg/server/telemetry"
	"github.com/cockroachdb/cockroach/pkg/sql/lex"
	"github.com/cockroachdb/cockroach/pkg/sql/sqltelemetry"
)

// AlterType represents an ALTER TYPE statement.
type AlterType struct {
	Type *UnresolvedObjectName
	Cmd  AlterTypeCmd
}

// Format implements the NodeFormatter interface.
func (node *AlterType) Format(ctx *FmtCtx) {
	ctx.WriteString("ALTER TYPE ")
	ctx.FormatNode(node.Type)
	ctx.FormatNode(node.Cmd)
}

// AlterTypeCmd represents a type modification operation.
type AlterTypeCmd interface {
	NodeFormatter
	alterTypeCmd()
	// TelemetryCounter returns the telemetry counter to increment
	// when this command is used.
	TelemetryCounter() telemetry.Counter
}

func (*AlterTypeAddValue) alterTypeCmd()    {}
func (*AlterTypeRenameValue) alterTypeCmd() {}
func (*AlterTypeRename) alterTypeCmd()      {}
func (*AlterTypeSetSchema) alterTypeCmd()   {}
func (*AlterTypeOwner) alterTypeCmd()       {}

var _ AlterTypeCmd = &AlterTypeAddValue{}
var _ AlterTypeCmd = &AlterTypeRenameValue{}
var _ AlterTypeCmd = &AlterTypeRename{}
var _ AlterTypeCmd = &AlterTypeSetSchema{}
var _ AlterTypeCmd = &AlterTypeOwner{}

// AlterTypeAddValue represents an ALTER TYPE ADD VALUE command.
type AlterTypeAddValue struct {
	NewVal      string
	IfNotExists bool
	Placement   *AlterTypeAddValuePlacement
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeAddValue) Format(ctx *FmtCtx) {
	ctx.WriteString(" ADD VALUE ")
	if node.IfNotExists {
		ctx.WriteString("IF NOT EXISTS ")
	}
	lex.EncodeSQLString(&ctx.Buffer, node.NewVal)
	if node.Placement != nil {
		if node.Placement.Before {
			ctx.WriteString(" BEFORE ")
		} else {
			ctx.WriteString(" AFTER ")
		}
		lex.EncodeSQLString(&ctx.Buffer, node.Placement.ExistingVal)
	}
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeAddValue) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "add_value")
}

// AlterTypeAddValuePlacement represents the placement clause for an ALTER
// TYPE ADD VALUE command ([BEFORE | AFTER] value).
type AlterTypeAddValuePlacement struct {
	Before      bool
	ExistingVal string
}

// AlterTypeRenameValue represents an ALTER TYPE RENAME VALUE command.
type AlterTypeRenameValue struct {
	OldVal string
	NewVal string
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeRenameValue) Format(ctx *FmtCtx) {
	ctx.WriteString(" RENAME VALUE ")
	lex.EncodeSQLString(&ctx.Buffer, node.OldVal)
	ctx.WriteString(" TO ")
	lex.EncodeSQLString(&ctx.Buffer, node.NewVal)
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeRenameValue) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "rename_value")
}

// AlterTypeRename represents an ALTER TYPE RENAME command.
type AlterTypeRename struct {
	NewName string
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeRename) Format(ctx *FmtCtx) {
	ctx.WriteString(" RENAME TO ")
	ctx.WriteString(node.NewName)
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeRename) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "rename")
}

// AlterTypeSetSchema represents an ALTER TYPE SET SCHEMA command.
type AlterTypeSetSchema struct {
	Schema string
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeSetSchema) Format(ctx *FmtCtx) {
	ctx.WriteString(" SET SCHEMA ")
	ctx.WriteString(node.Schema)
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeSetSchema) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "set_schema")
}

// AlterTypeOwner represents an ALTER TYPE OWNER TO command.
type AlterTypeOwner struct {
	Owner string
}

// Format implements the NodeFormatter interface.
func (node *AlterTypeOwner) Format(ctx *FmtCtx) {
	ctx.WriteString(" OWNER TO ")
	ctx.FormatNameP(&node.Owner)
}

// TelemetryCounter implements the AlterTypeCmd interface.
func (node *AlterTypeOwner) TelemetryCounter() telemetry.Counter {
	return sqltelemetry.SchemaChangeAlterCounterWithExtra("type", "owner")
}
