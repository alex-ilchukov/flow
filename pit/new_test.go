package pit_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/pit"
)

type miner[E errors.Senders] struct {
	total int
	last  int
	err   error
}

func (m *miner[E]) Mine(p pit.Pad[int, E]) {
	for ; m.last < m.total; m.last++ {
		err := p.Put(m.last)
		if err != nil {
			return
		}
	}

	if len(p.Errs()) > 0 && m.err != nil {
		p.Errs()[len(p.Errs())-1].Send(m.err)
	}
}

func TestResultWithNoErrors(t *testing.T) {
	ctx := context.Background()
	m := miner[errors.No]{total: 5}
	f := pit.New[int, errors.No](&m)
	err := flow.Run[int](ctx, f)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case m.last != 5:
		t.Errorf("got wrong last value: %d", m.last)
	}
}

func TestResultWithOneErrorWhenFlowIsSuccessful(t *testing.T) {
	ctx := context.Background()
	m := miner[errors.One]{total: 5}
	f := pit.New[int, errors.One](&m)
	err := flow.Run[int](ctx, f)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case m.last != 5:
		t.Errorf("got wrong last value: %d", m.last)
	}
}

func TestResultWithOneErrorWhenFlowIsErrorful(t *testing.T) {
	ctx := context.Background()
	err := fmt.Errorf("serious problem")
	m := miner[errors.One]{total: 5, err: err}
	f := pit.New[int, errors.One](&m)
	err = flow.Run[int](ctx, f)

	switch {
	case err == nil:
		t.Error("got no error")

	case err != m.err:
		t.Errorf("got invalid error: %v", err)

	case m.last != 5:
		t.Errorf("got wrong last value: %d", m.last)
	}
}
