package magnate

import (
	"fmt"
	"io"

	"github.com/rakyll/pb"
	"labix.org/v2/mgo"
)

type Namer interface {
	CollectionName() string
}

type Client struct {
	*mgo.Database
}

func (c Client) C(n Namer) *mgo.Collection {
	return c.Database.C(n.CollectionName())
}

type Operation interface {
	Describe() string
	Execute(Client) error
}

type Insert struct {
	Namer   Namer
	Content interface{}
}

func (i Insert) Describe() string {
	return fmt.Sprintf(
		"db.%s.insert(%#v)",
		i.Namer.CollectionName(),
		i.Content,
	)
}

func (i Insert) Execute(c Client) error {
	return c.C(i.Namer).Insert(i.Content)
}

type Update struct {
	Namer    Namer
	Selector interface{}
	Content  interface{}
}

func (u Update) Describe() string {
	return fmt.Sprintf(
		"db.%s.update(%#v, %#v)",
		u.Namer.CollectionName(),
		u.Selector,
		u.Content,
	)
}

func (u Update) Execute(c Client) error {
	return c.C(u.Namer).Update(u.Selector, u.Content)
}

type Remove struct {
	Namer    Namer
	Selector interface{}
}

func (r Remove) Describe() string {
	return fmt.Sprintf(
		"db.%s.remove(%#v)",
		r.Namer.CollectionName(),
		r.Selector,
	)
}

func (r Remove) Execute(c Client) error {
	return c.C(r.Namer).Remove(r.Selector)
}

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
