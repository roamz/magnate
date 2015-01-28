package magnate

import "testing"

func TestChangesStack(t *testing.T) {
	var ch Change
	cs := Changes{}

	cs.Push(ch)
	cs.Push(ch)

	if cs.Empty() {
		t.Error("should be two elements")
	}

	if cs.Pop() != ch {
		t.Error("should pop empty change")
	}

	if cs.Empty() {
		t.Error("should be one element")
	}

	if cs.Pop() != ch {
		t.Error("should pop empty change")
	}

	if !cs.Empty() {
		t.Error("should be empty")
	}

}
