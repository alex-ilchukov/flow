package chans_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/alex-ilchukov/flow/chans"
)

func TestDiscard(t *testing.T) {
	r := make([]int, 0)
	c := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		c <- 1
		c <- 2
		c <- 3
		r = append(r, 1)
		close(c)
		wg.Done()
	}()

	go func() {
		chans.Discard(c)
		r = append(r, 2)
		wg.Done()
	}()

	wg.Wait()
	if !reflect.DeepEqual(r, []int{1, 2}) {
		t.Errorf("wrong sequence r: %v", r)
	}
}
