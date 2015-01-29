package magnate

import (
	"errors"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// MigrationMarker attempts to get the Marker for a given Migration from the DB.
func MigrationMarker(c Client, m Migration) (
	marker Marker,
	status MigrationStatus,
	err error,
) {
	err = c.C(marker).Find(bson.M{"number": m.Number()}).One(&marker)

	if err == mgo.ErrNotFound {
		return marker, status, nil
	}

	status = Migrated
	if marker.Partial {
		status = Partial
	}

	return marker, status, err
}

type MigrationStatus int

const (
	Pending MigrationStatus = iota
	Partial
	Migrated
)

type Status struct {
	Migration
	MigrationStatus
}

func Statuses(
	c Client,
	migrations ...Migration,
) ([]Status, error) {
	var statuses []Status
	for _, migration := range migrations {
		_, status, err := MigrationMarker(c, migration)
		if err != nil {
			return statuses, checkError(err, migration)
		}

		statuses = append(statuses, Status{migration, status})
	}

	return statuses, nil
}

func CleanUp(n int, statuses ...Status) ([]Status, error) {
	return Clean(
		func(status Status) bool { return status.MigrationStatus == Migrated },
		n,
		statuses...,
	)
}

func CleanDown(n int, statuses ...Status) ([]Status, error) {
	return Clean(
		func(status Status) bool { return status.MigrationStatus == Pending },
		n,
		statuses...,
	)
}

func Clean(
	alreadyApplied func(Status) bool,
	number int,
	statuses ...Status,
) ([]Status, error) {
	var applying bool
	var cleaned []Status

	for _, status := range statuses {
		// Flip applying when we reach a migration that needs to be applied
		if !applying && !alreadyApplied(status) {
			applying = true
		}

		if applying && alreadyApplied(status) {
			err := errors.New("bad migration application")
			return cleaned, applicationError(err, status.Migration)
		}

		if applying {
			cleaned = append(cleaned, status)

			if number > 0 {
				if len(cleaned) >= number {
					return cleaned, nil
				}
			}
		}
	}

	return cleaned, nil
}

func (r Runner) MarkMigration(m Migration) error {
	if !r.NoDry {
		return nil
	}

	query := bson.M{"number": m.Number()}
	update := bson.M{"$unset": bson.M{"partial": true}}
	return r.FC.C(Marker{}).Update(query, update)
}

func (r Runner) MarkPartialMigration(status Status) error {
	if !r.NoDry {
		return nil
	}

	if status.MigrationStatus == Partial {
		return nil
	}

	marker := Marker{
		bson.NewObjectId(),
		status.Migration.Number(),
		status.Migration.Label(),
		true,
	}

	return r.FC.C(marker).Insert(marker)
}

func (r Runner) MarkDownMigration(m Migration) error {
	if !r.NoDry {
		return nil
	}

	return r.FC.C(Marker{}).Remove(bson.M{"number": m.Number()})
}

func (r Runner) MarkPartialDownMigration(status Status) error {
	if !r.NoDry {
		return nil
	}

	if status.MigrationStatus == Partial {
		return nil
	}

	query := bson.M{"number": status.Migration.Number()}
	update := bson.M{"$set": bson.M{"partial": true}}
	return r.FC.C(Marker{}).Update(query, update)
}
