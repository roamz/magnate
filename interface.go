package magnate

type Operation interface {
	Describe() string
	Execute(Client) error
}

// Change is an individual operation that can be performed on the database. It
// must have an up and down component so that it can be rolled forward or
// backwards.
type Change struct {
	Forwards Operation
	Rollback Operation
}

// Changes is an atomic set of changes that must be completed together.
type Changes []Change

type MaybeChanges struct {
	Changes Changes
	Err     error
}

type ChangeSet struct {
	Count int
	Func  func(chan<- MaybeChanges)
}

type Migration interface {
	Number() int
	Label() string

	Up(Client) (ChangeSet, error)
	Down(Client) (ChangeSet, error)
}
