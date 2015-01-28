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

/*
func (r Runner) Apply(changes Changes) error {
	for i, change := range changes {
		if r.Verbose {
			if _, cerr = fmt.Fprintln(r.Out, op.Describe()); err != nil {
				return err
			}
		}

	if r.NoDry {
		if err = op.Execute(r.Client); err != nil {
			return err
		}
	}
}


func (r Runner) Run(cs ChangeSet) error {
	var (
		bar *pb.ProgressBar
		err, cerr error
		mcc = make(chan MaybeChanges)
	)

	if r.ProgressBar {
		bar = pb.StartNew(cs.Count)
		bar.ShowSpeed = true
		defer bar.Finish()
	}

	go cs.Func(mcc)
	for changes := range mcc {
		if changes.Err != nil {
			err = changes.Err
		}

		if err != nil {
			continue
		}

		for _, change := range changes.Changes {
			if r.Verbose {
				if _, cerr = fmt.Fprintln(r.Out, op.Describe()); err != nil {
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
*/

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
