package emitters_test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/alex-ilchukov/flow/emitters"
	flowerrors "github.com/alex-ilchukov/flow/errors"
	"github.com/alex-ilchukov/flow/values"
)

type c[E flowerrors.Senders] struct {
	err error
}

func (c *c[E]) Convert(
	ctx  context.Context,
	r    values.Receiver[int],
	s    values.Sender[bool],
	errs E,
) {
	for {
		i, err := r.Receive()
		if err == values.Over {
			break
		}
		if err != nil {
			return
		}

		s.Send(i % 2 == 0)
	}

	if c.err != nil && len(errs) > 0 {
		errs[len(errs)-1].Send(c.err)
	}
}

func TestTransformWithNoErrorsEmitsSuccessfully(t *testing.T) {
	ctx := context.Background()
	p := p[flowerrors.No]{amount: 5}
	c := c[flowerrors.No]{}

	e0 := emitters.New[int, flowerrors.No](&p)
	e := emitters.Transform[int, bool, flowerrors.No](e0, &c)

	out, _ := e.Emit(ctx)
	a := []bool{}
	for b := range out {
		a = append(a, b)
	}

	if !reflect.DeepEqual(a, []bool{true, false, true, false, true}) {
		t.Errorf("Got wrong emitted values: %v", a)
	}
}

func TestTransformWithOneErrorEmitsSuccessfully(t *testing.T) {
	ctx := context.Background()
	p := p[flowerrors.No]{amount: 5}
	c := c[flowerrors.One]{}

	e0 := emitters.New[int, flowerrors.No](&p)
	e := emitters.Transform[int, bool, flowerrors.One](e0, &c)

	out, rerrs := e.Emit(ctx)
	rerr := flowerrors.Merge(ctx, rerrs)
	a := []bool{}

	err := error(nil)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for b := range out {
			a = append(a, b)
		}

		wg.Done()
	}()

	go func() {
		err = <-rerr
		wg.Done()
	}()

	wg.Wait()

	switch {
	case !reflect.DeepEqual(a, []bool{true, false, true, false, true}):
		t.Errorf("Got wrong emitted values: %v", a)

	case err != nil:
		t.Errorf("Got error %v", err)
	}
}

func TestTransformWithOneErrorEmitsWithConverterError(t *testing.T) {
	ctx := context.Background()
	p := p[flowerrors.No]{amount: 5}
	errConverterProblem := errors.New("convert has got serious problem")
	c := c[flowerrors.One]{err: errConverterProblem}

	e0 := emitters.New[int, flowerrors.No](&p)
	e := emitters.Transform[int, bool, flowerrors.One](e0, &c)

	out, rerrs := e.Emit(ctx)
	rerr := flowerrors.Merge(ctx, rerrs)
	a := []bool{}

	err := error(nil)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for b := range out {
			a = append(a, b)
		}

		wg.Done()
	}()

	go func() {
		err = <-rerr
		wg.Done()
	}()

	wg.Wait()

	switch {
	case !reflect.DeepEqual(a, []bool{true, false, true, false, true}):
		t.Errorf("Got wrong emitted values: %v", a)

	case err == nil:
		t.Errorf("Got no error")

	case err != errConverterProblem:
		t.Errorf("Got wrong error %v", err)
	}
}

func TestTransformWithOneErrorEmitsWithEmitterError(t *testing.T) {
	ctx := context.Background()
	errEmitterProblem := errors.New("emit has got serious problem")
	p := p[flowerrors.One]{amount: 5, err: errEmitterProblem}
	c := c[flowerrors.One]{}

	e0 := emitters.New[int, flowerrors.One](&p)
	e := emitters.Transform[int, bool, flowerrors.One](e0, &c)

	out, rerrs := e.Emit(ctx)
	rerr := flowerrors.Merge(ctx, rerrs)
	a := []bool{}

	err := error(nil)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for b := range out {
			a = append(a, b)
		}

		wg.Done()
	}()

	go func() {
		err = <-rerr
		wg.Done()
	}()

	wg.Wait()

	switch {
	case !reflect.DeepEqual(a, []bool{true, false, true, false, true}):
		t.Errorf("Got wrong emitted values: %v", a)

	case err == nil:
		t.Errorf("Got no error")

	case err != errEmitterProblem:
		t.Errorf("Got wrong error %v", err)
	}
}
