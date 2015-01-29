package main

import (
	"fmt"
	"os"
)

func (r Runner) Up(n int, migrations ...Migration) error {
	By(Ascending).Sort(migrations)

	statuses, err := Statuses(r.Client, migrations...)
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

		cs, err := status.Migration.Up(r.Client)
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
	By(Descending).Sort(migrations)

	statuses, err := Statuses(r.Client, migrations...)
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

		ops, err := status.Migration.Down(r.Client)
		if err != nil {
			return opGatherError(err, status.Migration)
		}

		if err = r.Run(ops...); err != nil {
			return opPerformError(err, status.Migration)
		}

		if err = r.UnMarkMigration(status.Migration); err != nil {
			return markError(err, status.Migration)
		}

	}

	return nil
}

func Status(r Runner) error {
	By(Ascending).Sort(migrations)

	for _, migration := range migrations {
		has, err := HasMigrated(r.Client, migration)
		if err != nil {
			return checkError(err, migration)
		}

		ch := ' '
		if has {
			ch = '*'
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
