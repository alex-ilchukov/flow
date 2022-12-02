package chans_test

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/alex-ilchukov/flow/chans"
)

func push(ch chan<- int, vals ...int) {
	for _, v := range vals {
		ch <- v
	}

	close(ch)
}

func TestMerge(t *testing.T) {
	ctx := context.Background()

	e1 := make(chan int)
	go push(e1, 1, 2, 3)

	e2 := make(chan int)
	go push(e2, 4, 5, 6)

	merged := chans.Merge(ctx, e1, e2)
	a := []int{}
	for v := range merged {
		a = append(a, v)
	}

	sort.Ints(a)
	if !reflect.DeepEqual(a, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("invalid values: %v", a)
	}
}

func TestMergeWhenCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	e1 := make(chan int)
	defer close(e1)

	e2 := make(chan int)
	defer close(e2)

	merged := chans.Merge(ctx, e1, e2)

	go cancel()

	a := []int{}
	for v := range merged {
		a = append(a, v)
	}

	if !reflect.DeepEqual(a, []int{}) {
		t.Errorf("invalid values: %v", a)
	}
}

func TestMergeWhenDeadlineExceeded(t *testing.T) {
	ctx := context.Background()
	deadline := time.Now().Add(time.Microsecond)
	ctx, _ = context.WithDeadline(ctx, deadline)

	e1 := make(chan int)
	defer close(e1)

	e2 := make(chan int)
	defer close(e2)

	merged := chans.Merge(ctx, e1, e2)

	a := []int{}
	for v := range merged {
		a = append(a, v)
	}

	if !reflect.DeepEqual(a, []int{}) {
		t.Errorf("invalid values: %v", a)
	}
}
