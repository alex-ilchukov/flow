package errors_test

import (
	"context"
	"errors"
	"testing"
	"time"

	flow_errors "github.com/alex-ilchukov/flow/errors"
)

func TestAnyWhenErrorPopped(t *testing.T) {
	ctx := context.Background()

	e1 := make(chan error)
	e2 := make(chan error)
	e3 := make(chan error)
	e4 := make(chan error)

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

	err := flow_errors.Any(
		ctx,
		[]<-chan error{e1, e2},
		[]<-chan error{e3, e4},
	)

	switch {
	case err == nil:
		t.Errorf("Error is not popped")

	case err != e:
		t.Errorf("Error is popped, but it is %v, not %v", err, e)
	}
}

func TestAnyWithNoError(t *testing.T) {
	ctx := context.Background()

	e1 := make(chan error)
	e2 := make(chan error)
	e3 := make(chan error)
	e4 := make(chan error)

	go close(e1)
	go close(e2)
	go close(e3)
	go close(e4)

	err := flow_errors.Any(
		ctx,
		[]<-chan error{e1, e2},
		[]<-chan error{e3, e4},
	)

	if err != nil {
		t.Errorf("Error %v is popped", err)
	}
}

func TestAnyWhenCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	e1 := make(chan error)
	e2 := make(chan error)
	e3 := make(chan error)
	e4 := make(chan error)

	go cancel()

	err := flow_errors.Any(
		ctx,
		[]<-chan error{e1, e2},
		[]<-chan error{e3, e4},
	)

	if err != nil {
		t.Errorf("Error %v is popped", err)
	}

	close(e1)
	close(e2)
	close(e3)
	close(e4)
}

func TestAnyWhenDeadlineExceeded(t *testing.T) {
	ctx := context.Background()
	deadline := time.Now().Add(time.Microsecond)
	ctx, _ = context.WithDeadline(ctx, deadline)

	e1 := make(chan error)
	e2 := make(chan error)
	e3 := make(chan error)
	e4 := make(chan error)

	err := flow_errors.Any(
		ctx,
		[]<-chan error{e1, e2},
		[]<-chan error{e3, e4},
	)

	if err != nil {
		t.Errorf("Error %v is popped", err)
	}

	close(e1)
	close(e2)
	close(e3)
	close(e4)
}
