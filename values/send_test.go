package values_test

import (
	"context"
	"testing"
	"time"

	"github.com/alex-ilchukov/flow/values"
)

func TestSendWhenSuccessful(t *testing.T) {
	ctx := context.Background()
	ch := make(chan int)
	v := 42
	vRead := 0

	go func() { vRead = <-ch }()

	status := values.Send(ctx, ch, v)

	switch {
	case vRead != v:
		t.Errorf("Invalid value %d read (should be %d)", vRead, v)

	case status != nil:
		t.Errorf("Got send error %v (should have no error)", status)
	}
}

const testSendWhenCanceledError = "Send hasn't been canceled: %v (should " +
	"have been canceled)"

func TestSendWhenCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan int)
	v := 42

	go cancel()

	status := values.Send(ctx, ch, v)

	if status != context.Canceled {
		t.Errorf(testSendWhenCanceledError, status)
	}
}

const testSendWhenDeadlineExceededError = "Send hasn't been exceeded " +
	"deadline: %v (should have been exceeded)"

func TestSendWhenDeadlineExceeded(t *testing.T) {
	ctx := context.Background()
	deadline := time.Now().Add(time.Microsecond)
	ctx, _ = context.WithDeadline(ctx, deadline)
	ch := make(chan int)
	v := 42
	status := values.Send(ctx, ch, v)

	if status != context.DeadlineExceeded {
		t.Errorf(testSendWhenDeadlineExceededError, status)
	}
}