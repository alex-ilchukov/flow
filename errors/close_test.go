package errors_test

import (
	"testing"

	"github.com/alex-ilchukov/flow/errors"
)

func TestClose(t *testing.T) {
	rerrs := make([]<-chan error, 1)
	werrs := make([]chan<- error, 1)
	for i := 0; i < len(werrs); i++ {
		ch := make(chan error)
		rerrs[i] := ch
		werrs[i] := ch
	}

	errors.Close(werrs)
	_, open := <-rerrs[0]
	if open {
		t.Errorf("Channel is still open after closing")
	}
}
