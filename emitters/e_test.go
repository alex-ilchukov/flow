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

type p[E flowerrors.Senders] struct {
	amount int
	err    error
}

func (p *p[E]) Produce(ctx context.Context, s values.Sender[int], errs E) {
	for i := 0; i < p.amount; i++ {
		s.Send(i)
	}

	if p.err != nil && len(errs) > 0 {
		errs[len(errs)-1].Send(p.err)
	}
}

func TestEmitterWithNoErrorsEmitsSuccessfully(t *testing.T) {
	ctx := context.Background()
	p := p[flowerrors.No]{amount: 5}
	e := emitters.New[int, flowerrors.No](&p)
	out, _ := e.Emit(ctx)
	a := []int{}
	for i := range out {
		a = append(a, i)
	}

	if !reflect.DeepEqual(a, []int{0, 1, 2, 3, 4}) {
		t.Errorf("Got wrong emitted values: %v", a)
	}
}

func TestEmitterWithOneErrorEmitsSuccessfully(t *testing.T) {
	ctx := context.Background()
	p := p[flowerrors.One]{amount: 5, err: nil}
	e := emitters.New[int, flowerrors.One](&p)
	out, rerrs := e.Emit(ctx)
	rerr := rerrs[0]
	a := []int{}
	err := error(nil)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for i := range out {
			a = append(a, i)
		}

		wg.Done()
	}()

	go func() {
		err = <-rerr
		wg.Done()
	}()

	wg.Wait()

	switch {
	case !reflect.DeepEqual(a, []int{0, 1, 2, 3, 4}):
		t.Errorf("Got wrong emitted values: %v", a)

	case err != nil:
		t.Errorf("Got error %v", err)
	}
}

func TestEmitterWithOneErrorEmitsWithError(t *testing.T) {
	ctx := context.Background()
	errProducerProblem := errors.New("producer has got serious problem")
	p := p[flowerrors.One]{amount: 5, err: errProducerProblem}
	e := emitters.New[int, flowerrors.One](&p)
	out, rerrs := e.Emit(ctx)
	rerr := rerrs[0]
	a := []int{}
	err := error(nil)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for i := range out {
			a = append(a, i)
		}

		wg.Done()
	}()

	go func() {
		err = <-rerr
		wg.Done()
	}()

	wg.Wait()

	switch {
	case !reflect.DeepEqual(a, []int{0, 1, 2, 3, 4}):
		t.Errorf("Got wrong emitted values: %v", a)

	case err == nil:
		t.Errorf("Got no error")

	case err != errProducerProblem:
		t.Errorf("Got wrong error %v", err)
	}
}
