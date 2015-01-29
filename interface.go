package magnate

// Operation is a
type Operation interface {
	Describe() string
	Execute(Client) error
}

// Change is an individual operation that can be performed on the database. It
// must have an up and down component so that it can be rolled forward or
// backwards.
type Change struct {
	Forwards Operation
	Revert   Operation
}

func (c Change) Reverse() Change {
	return Change{
		c.Revert,
		c.Forwards,
	}
}

// Changes is an atomic set of changes that must be completed together.
type Changes []Change

// FailingClient forms the first half of the migration handshake. A
// FailingClient is a client that performs the operations that a ChangeSet
// yields.
//
// If one of those operations fails, Err will be set accordingly. The migration
// should therefore watch for this err to be set, close the changes channel, and
// stop iterating if it is ever set.
//
// Note that this is purely an optimization to exit early, and if Err goes
// unchecked, the Runner will simply drain the channel without performing any
// further operations.
type FailingClient struct {
	Client
	Err error
}

func NewFailingClient(c Client) *FailingClient {
	return &FailingClient{Client: c}
}

// ChangeSet is a set of changes to be applied. Count indicates to the runner
// how many operations the set consists of. Err is an error that may be set at
// any point during the execution of Func, and will be checked by the Runner
// after the Changes channel has been closed.
type ChangeSet struct {
	Count int
	Func  func(chan<- Changes, *error)
}

type Migration interface {
	Number() int
	Label() string

	Up(*FailingClient) (*ChangeSet, error)
	Down(*FailingClient) (*ChangeSet, error)
}
