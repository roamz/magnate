package magnate

import (
	"errors"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func MigrationMarker(c Client, m Migration) (
	marker Marker,
	exists bool,
	err error,
) {
	err = c.C(marker).Find(bson.M{"number": m.Number()}).One(&marker)
	if err != mgo.ErrNotFound {
		exists = true
	}

	return marker, exists, err
}

type MigrationStatus struct {
	Migration
	Migrated bool
}

func Statuses(
	c Client,
	migrations ...Migration,
) ([]MigrationStatus, error) {
	var statuses []MigrationStatus
	for _, migration := range migrations {
		_, exists, err := MigrationMarker(c, migration)
		if err != nil {
			return statuses, err
			// TODO:FIX
			//return statuses, checkError(err, migration)
		}

		statuses = append(statuses, MigrationStatus{migration, exists})
	}

	return statuses, nil
}

func CleanUp(n int, statuses ...MigrationStatus) ([]MigrationStatus, error) {
	return Clean(
		func(status MigrationStatus) bool { return status.Migrated },
		n,
		statuses...,
	)
}

func CleanDown(n int, statuses ...MigrationStatus) ([]MigrationStatus, error) {
	return Clean(
		func(status MigrationStatus) bool { return !status.Migrated },
		n,
		statuses...,
	)
}

func Clean(
	alreadyApplied func(MigrationStatus) bool,
	number int,
	statuses ...MigrationStatus,
) ([]MigrationStatus, error) {
	var applying bool
	var cleaned []MigrationStatus

	for _, status := range statuses {
		// Flip applying when we reach a migration that needs to be applied
		if !applying && !alreadyApplied(status) {
			applying = true
		}

		if applying && alreadyApplied(status) {
			err := errors.New("bad migration application")
			//TODO:FIX
			//return cleaned, applicationError(err, status.Migration)
			return cleaned, err
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

func UnMarkMigration(r Runner, m Migration) error {
	if !r.NoDry {
		return nil
	}

	return r.Client.C(Marker{}).Remove(bson.M{"number": m.Number()})
}

func MarkMigration(r Runner, m Migration) error {
	if !r.NoDry {
		return nil
	}

	marker := Marker{
		bson.NewObjectId(),
		m.Number(),
		m.Label(),
		false,
	}

	return r.Client.C(marker).Insert(marker)
}
