package errors_test

import (
	"context"
	"errors"
	"testing"
	"time"

	flow_errors "github.com/alex-ilchukov/flow/errors"
)

func TestMergeWhenErrorPop(t *testing.T) {
	ctx := context.Background()

	e1 := make(chan error)
	e2 := make(chan error)
	e3 := make(chan error)
	e4 := make(chan error)

	merged := flow_errors.Merge(
		ctx,
		[]<-chan error{e1, e2},
		[]<-chan error{e3, e4},
	)

	e := errors.New("test")
	report := func(ch chan<- error, e error) {
		if e != nil {
			ch <- e
		}
		close(ch)
	}
	go report(e1, nil)
	go report(e2, nil)
	go report(e3, e)
	go report(e4, nil)

	p := error(nil)
	for p = range merged {
		if p != e {
			t.Errorf("Error is popped, but it is %v, not %v", p, e)
		}
	}

	if p == nil {
		t.Errorf("Error is not popped")
	}
}

func TestMergeWhenChannelsClosed(t *testing.T) {
	ctx := context.Background()

	e1 := make(chan error)
	e2 := make(chan error)
	e3 := make(chan error)
	e4 := make(chan error)

	go close(e1)
	go close(e2)
	go close(e3)
	go close(e4)

	merged := flow_errors.Merge(
		ctx,
		[]<-chan error{e1, e2},
		[]<-chan error{e3, e4},
	)

	for p := range merged {
		if p != nil {
			t.Errorf("Error %v is popped", p)
		}
	}
}

func TestMergeWhenCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	e1 := make(chan error)
	e2 := make(chan error)
	e3 := make(chan error)
	e4 := make(chan error)

	merged := flow_errors.Merge(
		ctx,
		[]<-chan error{e1, e2},
		[]<-chan error{e3, e4},
	)

	go cancel()

	for p := range merged {
		if p != nil {
			t.Errorf("Error is popped: %v", p)
		}
	}

	close(e1)
	close(e2)
	close(e3)
	close(e4)
}

func TestMergeWhenDeadlineExceeded(t *testing.T) {
	ctx := context.Background()
	deadline := time.Now().Add(time.Microsecond)
	ctx, _ = context.WithDeadline(ctx, deadline)

	e1 := make(chan error)
	e2 := make(chan error)
	e3 := make(chan error)
	e4 := make(chan error)

	merged := flow_errors.Merge(
		ctx,
		[]<-chan error{e1, e2},
		[]<-chan error{e3, e4},
	)

	for p := range merged {
		if p != nil {
			t.Errorf("Error is popped: %v", p)
		}
	}

	close(e1)
	close(e2)
	close(e3)
	close(e4)
}
