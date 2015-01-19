package magnate

type Operation interface {
	Describe() string
	Execute(Client) error
}

type Change struct {
	Up   Operation
	Down Operation
}

type Changes []Change

type Migration interface {
	Number() int
	Label() string
	Up(Client) (Changes, error)
	Down(Client) (Changes, error)
}
