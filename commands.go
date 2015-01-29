package magnate

import (
	"fmt"
	"os"
)

func (r Runner) Up(n int, migrations ...Migration) error {
	MigrationBy(Ascending).Sort(migrations)

	statuses, err := Statuses(r.FC.Client, migrations...)
	if err != nil {
		return err
	}

	statuses, err = CleanUp(n, statuses...)
	if err != nil {
		return err
	}

	for _, status := range statuses {
		fmt.Fprintf(
			os.Stderr,
			"[UP] %d:%s\n",
			status.Migration.Number(),
			status.Migration.Label(),
		)

		if err := r.MarkPartialMigration(status.Migration); err != nil {
			return markError(err, status.Migration)
		}

		cs, err := status.Migration.Up(r.FC)
		if err != nil {
			return opGatherError(err, status.Migration)
		}

		if err = r.Run(cs); err != nil {
			return opPerformError(err, status.Migration)
		}

		if err = r.MarkMigration(status.Migration); err != nil {
			return markError(err, status.Migration)
		}

	}

	return nil
}

func (r Runner) Down(n int, migrations ...Migration) error {
	MigrationBy(Descending).Sort(migrations)

	statuses, err := Statuses(r.FC.Client, migrations...)
	if err != nil {
		return err
	}

	statuses, err = CleanDown(n, statuses...)
	if err != nil {
		return err
	}

	for _, status := range statuses {
		fmt.Fprintf(
			os.Stderr,
			"[DOWN] %d:%s\n",
			status.Migration.Number(),
			status.Migration.Label(),
		)

		cs, err := status.Migration.Down(r.FC)
		if err != nil {
			return opGatherError(err, status.Migration)
		}

		if err = r.Run(cs); err != nil {
			return opPerformError(err, status.Migration)
		}

		if err = r.UnMarkMigration(status.Migration); err != nil {
			return markError(err, status.Migration)
		}

	}

	return nil
}

func (r Runner) Status(migrations ...Migration) error {
	MigrationBy(Ascending).Sort(migrations)

	for _, migration := range migrations {
		_, status, err := MigrationMarker(r.FC.Client, migration)
		if err != nil {
			return checkError(err, migration)
		}

		var ch rune
		switch status {
		case Pending:
			ch = ' '
		case Partial:
			ch = '-'
		case Migrated:
			ch = '*'
		default:
			ch = '?'
		}

		fmt.Fprintf(
			os.Stderr,
			"[%c] %d:%s\n",
			ch, migration.Number(),
			migration.Label(),
		)
	}

	return nil
}
