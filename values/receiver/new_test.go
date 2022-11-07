package receiver_test

import (
	"context"
	"testing"
	"time"

	"github.com/alex-ilchukov/flow/values"
	"github.com/alex-ilchukov/flow/values/receiver"
)

func TestReceiveWhenSuccessful(t *testing.T) {
	ctx := context.Background()
	ch := make(chan int)
	r := receiver.New(ctx, ch)
	v := 42

	go func() { ch <- v }()

	u, err := r.Receive()

	switch {
	case u != v:
		t.Errorf("Invalid value %d is read (should be %d)", u, v)

	case err != nil:
		t.Errorf("Got receive error %v (should have no error)", err)
	}
}

const testReceiveWhenCanceledError = "Receive hasn't been canceled: %v " +
	"(should have been canceled)"

func TestReceiveWhenCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan int)
	r := receiver.New(ctx, ch)

	go cancel()

	_, err := r.Receive()

	if err != context.Canceled {
		t.Errorf(testReceiveWhenCanceledError, err)
	}
}

const testReceiveWhenClosedError = "Receive got invalid error on closed " +
	"channel: %v (should have been %v)"

func TestReceiveWhenClosed(t *testing.T) {
	ctx := context.Background()
	ch := make(chan int)
	r := receiver.New(ctx, ch)
	close(ch)

	_, err := r.Receive()

	if err != values.Over {
		t.Errorf(testReceiveWhenClosedError, err, values.Over)
	}
}

const testReceiveWhenDeadlineExceededError = "Receive hasn't been exceeded " +
	"deadline: %v (should have been exceeded)"

func TestReceiveWhenDeadlineExceeded(t *testing.T) {
	ctx := context.Background()
	deadline := time.Now().Add(time.Microsecond)
	ctx, _ = context.WithDeadline(ctx, deadline)
	ch := make(chan int)
	r := receiver.New(ctx, ch)
	_, err := r.Receive()

	if err != context.DeadlineExceeded {
		t.Errorf(testReceiveWhenDeadlineExceededError, err)
	}
}
