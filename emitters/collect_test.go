package emitters_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/alex-ilchukov/flow/emitters"
	flowerrors "github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values"
)

type con[E flowerrors.Senders] struct {
	ints []int
}

var errOverflow = errors.New("overflow")

func (c *con[E]) Consume(ctx context.Context, r values.Receiver[int], errs E) {
	for {
		i, status := r.Receive()
		if status != nil {
			return
		}

		if len(errs) > 0 && len(c.ints) >= cap(c.ints) {
			errs[len(errs)-1].Send(errOverflow)
			return
		}

		c.ints = append(c.ints, i)
	}
}

func TestCollectWithNoErrorsCollectsSuccessfully(t *testing.T) {
	ctx := context.Background()
	p := p[flowerrors.No]{amount: 5}
	c := con[flowerrors.No]{}

	e0 := emitters.New[int, flowerrors.No](&p)
	e := emitters.Collect[int, flowerrors.No](e0, &c)

	out, _ := e.Emit(ctx)
	<-out

	if !reflect.DeepEqual(c.ints, []int{0, 1, 2, 3, 4}) {
		t.Errorf("Got wrong collected values: %v", c.ints)
	}
}

func TestCollectWithOneErrorCollectsSuccessfully(t *testing.T) {
	ctx := context.Background()
	p := p[flowerrors.No]{amount: 5}
	c := con[flowerrors.One]{ints: make([]int, 0, 5)}

	e0 := emitters.New[int, flowerrors.No](&p)
	e := emitters.Collect[int, flowerrors.One](e0, &c)

	out, rerrs := e.Emit(ctx)
	rerr := rerrs[0]
	err := error(nil)

	select {
	case <-out:
		// nothing to do

	case err = <-rerr:
		// nothing to do
	}

	switch {
	case !reflect.DeepEqual(c.ints, []int{0, 1, 2, 3, 4}):
		t.Errorf("Got wrong collected values: %v", c.ints)

	case err != nil:
		t.Errorf("Got error %v", err)
	}
}

func TestCollectWithOneErrorCollectsWithError(t *testing.T) {
	ctx := context.Background()
	p := p[flowerrors.No]{amount: 5}
	c := con[flowerrors.One]{ints: make([]int, 0, 4)}

	e0 := emitters.New[int, flowerrors.No](&p)
	e := emitters.Collect[int, flowerrors.One](e0, &c)

	out, rerrs := e.Emit(ctx)
	rerr := rerrs[0]
	err := error(nil)

	select {
	case <-out:
		// nothing to do

	case err = <-rerr:
		// nothing to do
	}

	switch {
	case !reflect.DeepEqual(c.ints, []int{0, 1, 2, 3}):
		t.Errorf("Got wrong collected values: %v", c.ints)

	case err == nil:
		t.Errorf("Got no error")

	case err != errOverflow:
		t.Errorf("Got wrong error: %v", err)
	}
}
