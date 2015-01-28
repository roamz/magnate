package magnate

import (
	"fmt"
	"io"

	"github.com/rakyll/pb"
)

type OpErr struct {
	error
	Operation
}

func Describe(out io.Writer, ops ...Operation) error {
	var err error
	for _, op := range ops {
		if _, err = fmt.Fprintln(out, op.Describe()); err != nil {
			return err
		}
	}

	return err
}

func Execute(c Client, ops ...Operation) error {
	var err error
	for _, op := range ops {
		if err = op.Execute(c); err != nil {
			return OpErr{err, op}
		}
	}

	return err
}

type Runner struct {
	Client
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
		if err = op.Execute(r.Client); err != nil {
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

func (r Runner) Run(cs ChangeSet) error {
	var (
		bar *pb.ProgressBar
		err error
		mcc = make(chan MaybeChanges)
	)

	if r.ProgressBar {
		bar = pb.StartNew(cs.Count)
		bar.ShowSpeed = true
		defer bar.Finish()
	}

	go cs.Func(mcc)
	for changes := range mcc {
		if err = r.Apply(changes.Changes); err != nil {
			// signal mcc to exit
			break
		}
	}

	return err
}

func (r Runner) Run(ops ...Operation) error {
	var bar *pb.ProgressBar
	if r.ProgressBar {
		bar = pb.StartNew(len(ops))
		bar.ShowSpeed = true
		defer bar.Finish()
	}

	var err error

	for _, op := range ops {
		if r.ProgressBar {
			bar.Increment()
		}

		if r.Verbose {
			if _, err = fmt.Fprintln(r.Out, op.Describe()); err != nil {
				return err
			}
		}

		if r.NoDry {
			if err = op.Execute(r.Client); err != nil {
				return err
			}
		}
	}

	return err
}
