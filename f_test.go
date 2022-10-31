package flow_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/alex-ilchukov/flow"
	"github.com/alex-ilchukov/flow/values"
)

type e struct {
	amount int
	err    error
	sleep  time.Duration
}

var errEmitterProblem = errors.New("emitter has got real problem")

func (e *e) Emit(ctx context.Context) (<-chan int, []<-chan error) {
	out := make(chan int)
	err := make(chan error)

	go e.process(ctx, out, err)

	return out, []<-chan error{err}
}

func (e *e) process(ctx context.Context, out chan int, err chan error) {
	defer close(out)
	defer close(err)

	for i := 0; i < e.amount; i++ {
		if values.Send(ctx, out, i) != nil {
			return
		}

		time.Sleep(e.sleep)
	}

	if e.err != nil {
		values.Send(ctx, err, e.err)
	}
}

type c struct {
	ints  []int
	limit int
}

func (c *c) Collect(ctx context.Context, in <-chan int) []<-chan error {
	err := make(chan error)

	go c.process(ctx, in, err)

	return []<-chan error{err}
}

var errOverflow = errors.New("overflow")

func (c *c) process(ctx context.Context, in <-chan int, err chan error) {
	defer close(err)

	for {
		i, status := values.Receive(ctx, in)
		if status != nil {
			return
		}

		if len(c.ints) >= c.limit {
			values.Send(ctx, err, errOverflow)
			return
		}

		c.ints = append(c.ints, i)
	}
}

func TestFEmitter(t *testing.T) {
	e := &e{amount: 5}
	c := &c{limit: 10}
	f := flow.New[int](e, c)

	if f.Emitter() == nil {
		t.Errorf("Emitter is nil")
	}
}

func TestFCollector(t *testing.T) {
	e := &e{amount: 5}
	c := &c{limit: 10}
	f := flow.New[int](e, c)

	if f.Collector() == nil {
		t.Errorf("Collector is nil")
	}
}

func TestFRunWithNoError(t *testing.T) {
	ctx := context.Background()
	e := &e{amount: 5}
	c := &c{limit: 10}
	f := flow.New[int](e, c)
	err := f.Run(ctx)

	switch {
	case err != nil:
		t.Errorf("Error %v has appeared", err)

	case !reflect.DeepEqual(c.ints, []int{0, 1, 2, 3, 4}):
		t.Errorf("Collector hasn't gotten all the ints: %v", c.ints)
	}
}

func TestFRunWithEmitterError(t *testing.T) {
	ctx := context.Background()
	e := &e{amount: 5, err: errEmitterProblem}
	c := &c{limit: 10}
	f := flow.New[int](e, c)
	err := f.Run(ctx)

	if err != errEmitterProblem {
		t.Errorf("Error %v is not the %v", err, errEmitterProblem)
	}
}

func TestFRunWithCollectorError(t *testing.T) {
	ctx := context.Background()
	e := &e{amount: 5}
	c := &c{limit: 4}
	f := flow.New[int](e, c)
	err := f.Run(ctx)

	if err != errOverflow {
		t.Errorf("Error %v is not the %v", err, errOverflow)
	}
}

func TestFRunWhenCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	e := &e{amount: 5, sleep: time.Second}
	c := &c{limit: 10}
	f := flow.New[int](e, c)

	go cancel()

	err := f.Run(ctx)
	if err != nil {
		t.Errorf("Error %v has appeared", err)
	}
}

func TestFRunWhenDeadlineExceeded(t *testing.T) {
	ctx := context.Background()
	deadline := time.Now().Add(time.Microsecond)
	ctx, _ = context.WithDeadline(ctx, deadline)

	e := &e{amount: 5, sleep: time.Second}
	c := &c{limit: 10}
	f := flow.New[int](e, c)

	err := f.Run(ctx)
	if err != nil {
		t.Errorf("Error %v has appeared", err)
	}
}
