package stage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/stage"
	"github.com/alex-ilchukov/flow/values"
)

type fimpl struct {
	total int
	err   error
}

func (f *fimpl) Flow(ctx context.Context) (<-chan int, []<-chan error) {
	c := make(chan int)
	e := make(chan error)
	go f.process(ctx, c, e)

	return c, []<-chan error{e}
}

func (f *fimpl) process(ctx context.Context, c chan int, e chan error) {
	defer close(e)
	defer close(c)

	for i := 0; i < f.total; i++ {
		values.Send(ctx, c, i)
	}

	if f.err != nil {
		values.Send(ctx, e, f.err)
	}
}

type former struct {
	last int
	err  error
}

func (f *former) Form(j stage.Joint[int, int]) {

loop:
	for {
		v, err := j.Get()
		switch {
		case err == values.Over:
			break loop

		case err != nil:
			return
		}

		w := v * v
		err = j.Put(w)
		if err != nil {
			return
		}

		f.last = w
	}

	if f.err != nil {
		j.Report(f.err)
	}
}

type miner struct {
	total int
	last  int
	err   error
}

func (m *miner) Form(j stage.Joint[int, int]) {
	for ; m.last < m.total; m.last++ {
		err := j.Put(m.last)
		if err != nil {
			return
		}
	}

	if m.err != nil {
		j.Report(m.err)
	}
}

func TestResultWhenSuccessful(t *testing.T) {
	ctx := context.Background()
	f := fimpl{total: 5}
	former := former{}
	newf := stage.New[int, int](&f, &former)
	err := flow.Run[int](ctx, newf)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case former.last != 16:
		t.Errorf("got wrong last value: %d", former.last)
	}
}

func TestResultWhenSuccessfulAndFlowIsNil(t *testing.T) {
	ctx := context.Background()
	former := miner{total: 5}
	newf := stage.New[int, int](nil, &former)
	err := flow.Run[int](ctx, newf)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case former.last != 5:
		t.Errorf("got wrong last value: %d", former.last)
	}
}

func TestResultWhenFlowIsErrorful(t *testing.T) {
	ctx := context.Background()
	err := fmt.Errorf("serious problem")
	f := fimpl{total: 5, err: err}
	former := former{}
	newf := stage.New[int, int](&f, &former)
	err = flow.Run[int](ctx, newf)

	switch {
	case err == nil:
		t.Error("got no error")

	case err != f.err:
		t.Errorf("got invalid error: %v", err)

	case former.last != 16:
		t.Errorf("got wrong last value: %d", former.last)
	}
}

func TestResultWhenFormerIsErrorful(t *testing.T) {
	ctx := context.Background()
	f := fimpl{total: 5}
	err := fmt.Errorf("serious problem")
	former := former{err: err}
	newf := stage.New[int, int](&f, &former)
	err = flow.Run[int](ctx, newf)

	switch {
	case err == nil:
		t.Error("got no error")

	case err != former.err:
		t.Errorf("got invalid error: %v", err)

	case former.last != 16:
		t.Errorf("got wrong last value: %d", former.last)
	}
}

func TestResultWhenFormerIsErrorfulAndFlowIsNil(t *testing.T) {
	ctx := context.Background()
	err := fmt.Errorf("serious problem")
	former := miner{total: 5, err: err}
	newf := stage.New[int, int](nil, &former)
	err = flow.Run[int](ctx, newf)

	switch {
	case err == nil:
		t.Error("got no error")

	case err != former.err:
		t.Errorf("got invalid error: %v", err)

	case former.last != 5:
		t.Errorf("got wrong last value: %d", former.last)
	}
}
