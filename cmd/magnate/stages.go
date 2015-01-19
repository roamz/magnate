package migrations

import (
	"fmt"

	"github.com/dagoof/magnate"
)

// Stages are various places that might return an error during migration
type Stage int

const (
	MigrationCheck Stage = iota
	OrderCheck
	OpGather
	OpPerform
	MigrationMark
)

type MigrationError struct {
	error
	Migration
	Stage
	Description string
}

func (e MigrationError) Error() string {
	return e.Description
}

func checkError(err error, migration Migration) MigrationError {
	return MigrationError{
		err,
		migration,
		MigrationCheck,
		fmt.Sprintf(
			"failed to check if migration %d exists",
			migration.Number(),
		),
	}
}

func applicationError(err error, migration Migration) MigrationError {
	return MigrationError{
		err,
		migration,
		OrderCheck,
		fmt.Sprintf(
			"migration %d:%s has been incorrectly applied and must be reverted",
			migration.Number(),
			migration.Label(),
		),
	}
}

func opGatherError(err error, migration Migration) MigrationError {
	return MigrationError{
		err,
		migration,
		OpGather,
		fmt.Sprintf(
			"failed to gather operations for migration %d:%s",
			migration.Number(),
			migration.Label(),
		),
	}
}

func opPerformError(err error, migration Migration) MigrationError {
	var description string
	switch v := err.(type) {
	case magnate.OpErr:
		description = fmt.Sprintf(
			"migration %d:%s failed on operation %s",
			migration.Number(),
			migration.Label(),
			v.Operation.Describe(),
		)
	default:
		description = fmt.Sprintf(
			"migration %d:%s failed during operation",
			migration.Number(),
			migration.Label(),
		)
	}

	return MigrationError{
		err,
		migration,
		OpPerform,
		description,
	}
}

func markError(err error, migration Migration) MigrationError {
	return MigrationError{
		err,
		migration,
		MigrationMark,
		fmt.Sprintf(
			"failed to mark completed migration %d:%s",
			migration.Number(),
			migration.Label(),
		),
	}
}
