package magnate

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

type Migration struct {
	Client
	Verbose bool
	NoDry   bool
}

func (m Migration) Perform(ops ...Operation) error {
	if m.Verbose {
		for _, op := range ops {
			log.Println(op.Describe())
		}
	}

	if m.NoDry {
		for _, op := range ops {
			if err := op.Execute(m.Client); err != nil {
				return err
			}
		}
	}

	return nil
}