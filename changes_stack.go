package magnate

func (cs Changes) Empty() bool    { return len(cs) == 0 }
func (cs Changes) Peek() Change   { return cs[len(cs)-1] }
func (cs *Changes) Push(c Change) { *cs = append(*cs, c) }
func (cs *Changes) Pop() Change {
	c := cs.Peek()
	*cs = (*cs)[:len(*cs)-1]
	return c
}
