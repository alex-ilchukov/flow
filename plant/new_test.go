package plant_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/plant"
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

type former[E errors.Senders] struct {
	last int
	err  error
}

func (f *former[E]) Form(j plant.Joint[int, int, E]) {

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

	if len(j.Errs()) > 0 && f.err != nil {
		j.Errs()[len(j.Errs())-1].Send(f.err)
	}
}

type miner[E errors.Senders] struct {
	total int
	last  int
	err   error
}

func (m *miner[E]) Form(j plant.Joint[int, int, E]) {
	for ; m.last < m.total; m.last++ {
		err := j.Put(m.last)
		if err != nil {
			return
		}
	}

	if len(j.Errs()) > 0 && m.err != nil {
		j.Errs()[len(j.Errs())-1].Send(m.err)
	}
}

func TestResultWithNoErrors(t *testing.T) {
	ctx := context.Background()
	f := fimpl{total: 5}
	former := former[errors.No]{}
	newf := plant.New[int, int, errors.No](&f, &former)
	err := flow.Run[int](ctx, newf)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case former.last != 16:
		t.Errorf("got wrong last value: %d", former.last)
	}
}

func TestResultWithNoErrorsWhenFlowIsNil(t *testing.T) {
	ctx := context.Background()
	former := miner[errors.No]{total: 5}
	newf := plant.New[int, int, errors.No](nil, &former)
	err := flow.Run[int](ctx, newf)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case former.last != 5:
		t.Errorf("got wrong last value: %d", former.last)
	}
}

func TestResultWithOneErrorWhenSuccessful(t *testing.T) {
	ctx := context.Background()
	f := fimpl{total: 5}
	former := former[errors.One]{}
	newf := plant.New[int, int, errors.One](&f, &former)
	err := flow.Run[int](ctx, newf)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case former.last != 16:
		t.Errorf("got wrong last value: %d", former.last)
	}
}

func TestResultWithOneErrorWhenSuccessfulAndFlowIsNil(t *testing.T) {
	ctx := context.Background()
	former := miner[errors.One]{total: 5}
	newf := plant.New[int, int, errors.One](nil, &former)
	err := flow.Run[int](ctx, newf)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case former.last != 5:
		t.Errorf("got wrong last value: %d", former.last)
	}
}

func TestResultWithOneErrorWhenFlowIsErrorful(t *testing.T) {
	ctx := context.Background()
	err := fmt.Errorf("serious problem")
	f := fimpl{total: 5, err: err}
	former := former[errors.One]{}
	newf := plant.New[int, int, errors.One](&f, &former)
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

func TestResultWithOneErrorWhenFormerIsErrorful(t *testing.T) {
	ctx := context.Background()
	f := fimpl{total: 5}
	err := fmt.Errorf("serious problem")
	former := former[errors.One]{err: err}
	newf := plant.New[int, int, errors.One](&f, &former)
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

func TestResultWithOneErrorWhenFormerIsErrorfulAndFlowIsNil(t *testing.T) {
	ctx := context.Background()
	err := fmt.Errorf("serious problem")
	former := miner[errors.One]{total: 5, err: err}
	newf := plant.New[int, int, errors.One](nil, &former)
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
