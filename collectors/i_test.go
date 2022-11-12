package collectors_test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/alex-ilchukov/flow/collectors"
	flowerrors "github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values"
)

type c[E flowerrors.Senders] struct {
	ints []int
}

var errOverflow = errors.New("overflow")

func (c *c[E]) Consume(ctx context.Context, in <-chan int, errs E) {
	for {
		i, status := values.Receive(ctx, in)
		if status != nil {
			return
		}

		if len(errs) > 0 && len(c.ints) >= cap(c.ints) {
			values.Send(ctx, errs[len(errs)-1], errOverflow)
			return
		}

		c.ints = append(c.ints, i)
	}
}

func TestCollectorWithNoErrorsCollectsSuccessfully(t *testing.T) {
	ctx := context.Background()
	c := c[flowerrors.No]{}
	col := collectors.New[int, flowerrors.No](&c)

	ch := make(chan int)
	col.Collect(ctx, ch)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)

	if !reflect.DeepEqual(c.ints, []int{0, 1, 2, 3, 4}) {
		t.Errorf("Got wrong collected values: %v", c.ints)
	}
}

func TestCollectorWithOneErrorCollectsSuccessfully(t *testing.T) {
	ctx := context.Background()
	c := c[flowerrors.One]{ints: make([]int, 0, 5)}
	col := collectors.New[int, flowerrors.One](&c)
	ch := make(chan int)
	rerrs := col.Collect(ctx, ch)
	rerr := rerrs[0]
	err := error(nil)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)

		wg.Done()
	}()

	go func() {
		err = <-rerr
		wg.Done()
	}()

	wg.Wait()

	switch {
	case !reflect.DeepEqual(c.ints, []int{0, 1, 2, 3, 4}):
		t.Errorf("Got wrong collected values: %v", c.ints)

	case err != nil:
		t.Errorf("Got error %v", err)
	}
}

func TestCollectorWithOneErrorCollectsWithError(t *testing.T) {
	ctx := context.Background()
	c := c[flowerrors.One]{ints: make([]int, 0, 4)}
	col := collectors.New[int, flowerrors.One](&c)
	ch := make(chan int)
	rerrs := col.Collect(ctx, ch)
	rerr := rerrs[0]
	err := error(nil)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch)

		wg.Done()
	}()

	go func() {
		err = <-rerr
		wg.Done()
	}()

	wg.Wait()

	switch {
	case !reflect.DeepEqual(c.ints, []int{0, 1, 2, 3}):
		t.Errorf("Got wrong collected values: %v", c.ints)

	case err == nil:
		t.Errorf("Got no error")

	case err != errOverflow:
		t.Errorf("Got wrong error: %v", err)
	}
}
