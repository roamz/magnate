package magnate

import (
	"fmt"
	"io"

	"github.com/cheggaaa/pb"
)

type OpErr struct {
	error
	Operation
}

type Runner struct {
	FC          *FailingClient
	Out         io.Writer
	Verbose     bool
	NoDry       bool
	ProgressBar bool
}

func (r Runner) Something(op Operation) error {
	var err error

	if r.Verbose {
		if _, err = fmt.Fprintln(r.Out, op.Describe()); err != nil {
			return err
		}
	}

	if r.NoDry {
		if err = op.Execute(r.FC.Client); err != nil {
			return err
		}
	}

	return err
}

type RevertError struct {
	error        // implicit forwards error
	Revert error // extra revert error
}

func (r Runner) Apply(forwards Changes) error {
	var (
		change    Change
		revert    Changes
		err, rerr error
	)

	for _, change = range forwards {
		if err = r.Something(change.Forwards); err != nil {
			break
		}

		revert.Push(change)
	}

	if err == nil {
		return err
	}

	for !revert.Empty() {
		change = revert.Pop()
		if rerr = r.Something(change.Revert); rerr != nil {
			return RevertError{err, rerr}
		}
	}

	return err
}

func (r Runner) Run(cs *ChangeSet) error {
	var (
		bar *pb.ProgressBar
		err error
		mcc = make(chan Changes)
	)

	if r.ProgressBar {
		bar = pb.StartNew(cs.Count)
		bar.ShowSpeed = true
		defer bar.Finish()
	}

	go cs.Func(mcc)
	for changes := range mcc {
		if r.FC.Err != nil {
			continue
		}

		if err = r.Apply(changes); err != nil {
			r.FC.Err = err
			continue
		}

		if r.ProgressBar {
			bar.Increment()
		}
	}

	return err
}
