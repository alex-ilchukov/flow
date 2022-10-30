package values_test

import (
	"context"
	"testing"
	"time"

	"github.com/alex-ilchukov/flow/values"
)

func TestReceiveWhenSuccessful(t *testing.T) {
	ctx := context.Background()
	ch := make(chan int)
	v := 42

	go func() { ch <- v }()

	u, status := values.Receive(ctx, ch)

	switch {
	case u != v:
		t.Errorf("Invalid value %d is read (should be %d)", u, v)

	case status != nil:
		t.Errorf("Got receive error %v (should have no error)", status)
	}
}

const testReceiveWhenCanceledError = "Receive hasn't been canceled: %v " +
	"(should have been canceled)"

func TestReceiveWhenCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan int)

	go cancel()

	_, status := values.Receive(ctx, ch)

	if status != context.Canceled {
		t.Errorf(testReceiveWhenCanceledError, status)
	}
}

const testReceiveWhenClosedError = "Receive got invalid error on closed " +
	"channel: %v (should have been %v)"

func TestReceiveWhenClosed(t *testing.T) {
	ctx := context.Background()
	ch := make(chan int)
	close(ch)

	_, status := values.Receive(ctx, ch)

	if status != values.ErrClosedChannel {
		t.Errorf(
			testReceiveWhenClosedError,
			status,
			values.ErrClosedChannel,
		)
	}
}

const testReceiveWhenDeadlineExceededError = "Receive hasn't been exceeded " +
	"deadline: %v (should have been exceeded)"

func TestReceiveWhenDeadlineExceeded(t *testing.T) {
	ctx := context.Background()
	deadline := time.Now().Add(time.Microsecond)
	ctx, _ = context.WithDeadline(ctx, deadline)
	ch := make(chan int)
	_, status := values.Receive(ctx, ch)

	if status != context.DeadlineExceeded {
		t.Errorf(testReceiveWhenDeadlineExceededError, status)
	}
}
