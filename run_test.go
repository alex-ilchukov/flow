package flow_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alex-ilchukov/flow"
)

type noerrs struct {
	total   int
	written int
}

func (f *noerrs) Flow(context.Context) (<-chan int, []<-chan error) {
	c := make(chan int)
	go f.process(c)

	return c, nil
}

func (f *noerrs) process(c chan int) {
	defer close(c)

	for ; f.written < f.total; f.written++ {
		c <- f.written
	}
}

type witherrs struct {
	total    int
	err      error
	written  int
	canceled bool
}

func (f *witherrs) Flow(ctx context.Context) (<-chan int, []<-chan error) {
	c := make(chan int)
	e := make(chan error)
	go f.process(ctx, c, e)

	return c, []<-chan error{e}
}

func (f *witherrs) process(ctx context.Context, c chan int, e chan error) {
	defer close(e)
	defer close(c)

	if f.total < 0 {
		for !f.canceled {
			f.push(ctx, c, 0)
		}
	} else {
		for ; f.written < f.total && !f.canceled; f.written++ {
			f.push(ctx, c, f.written)
		}
	}

	if !f.canceled && f.err != nil {
		e <- f.err
	}
}

func (f *witherrs) push(ctx context.Context, c chan int, v int) {
	select {
	case <-ctx.Done():
		f.canceled = true

	case c <- f.written:
		// do nothing
	}
}

func TestRunWhenFlowIsNil(t *testing.T) {
	ctx := context.Background()
	err := flow.Run[int](ctx, nil)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestRunWhenFlowHasNoErrs(t *testing.T) {
	ctx := context.Background()
	f := noerrs{total: 5}
	err := flow.Run[int](ctx, &f)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case f.written != 5:
		t.Errorf("invalid write tries amount: %d", f.written)
	}
}

func TestRunWhenFlowIsErrorless(t *testing.T) {
	ctx := context.Background()
	f := witherrs{total: 5}
	err := flow.Run[int](ctx, &f)

	switch {
	case err != nil:
		t.Errorf("got error: %v", err)

	case f.written != 5:
		t.Errorf("invalid write tries amount: %d", f.written)
	}
}

func TestRunWhenFlowIsErrorful(t *testing.T) {
	ctx := context.Background()
	flowerr := errors.New("serious problem")
	f := witherrs{total: 5, err: flowerr}
	err := flow.Run[int](ctx, &f)

	switch {
	case err == nil:
		t.Error("got no error")

	case err != flowerr:
		t.Errorf("got invalid error: %v", err)

	case f.written != 5:
		t.Errorf("invalid write tries amount: %d", f.written)
	}
}

func TestRunWhenFlowIsCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go cancel()

	f := witherrs{total: -1}
	err := flow.Run[int](ctx, &f)

	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestRunWhenFlowHasReachedDeadline(t *testing.T) {
	ctx := context.Background()
	deadline := time.Now().Add(1 * time.Microsecond)
	ctx, _ = context.WithDeadline(ctx, deadline)

	f := witherrs{total: -1}
	err := flow.Run[int](ctx, &f)

	if err != nil {
		t.Errorf("got error: %v", err)
	}
}
