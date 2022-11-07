package sender_test

import (
	"context"
	"testing"
	"time"

	"github.com/alex-ilchukov/flow/values/sender"
)

func TestSendWhenSuccessful(t *testing.T) {
	ctx := context.Background()
	ch := make(chan int)
	s := sender.New(ctx, ch)
	v := 42
	vRead := 0

	go func() { vRead = <-ch }()

	err := s.Send(v)

	switch {
	case vRead != v:
		t.Errorf("Invalid value %d read (should be %d)", vRead, v)

	case err != nil:
		t.Errorf("Got send error %v (should have no error)", err)
	}
}

const testSendWhenCanceledError = "Send hasn't been canceled: %v (should " +
	"have been canceled)"

func TestSendWhenCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan int)
	s := sender.New(ctx, ch)
	v := 42

	go cancel()

	err := s.Send(v)

	if err != context.Canceled {
		t.Errorf(testSendWhenCanceledError, err)
	}
}

const testSendWhenDeadlineExceededError = "Send hasn't been exceeded " +
	"deadline: %v (should have been exceeded)"

func TestSendWhenDeadlineExceeded(t *testing.T) {
	ctx := context.Background()
	deadline := time.Now().Add(time.Microsecond)
	ctx, _ = context.WithDeadline(ctx, deadline)
	ch := make(chan int)
	s := sender.New(ctx, ch)
	v := 42
	err := s.Send(v)

	if err != context.DeadlineExceeded {
		t.Errorf(testSendWhenDeadlineExceededError, err)
	}
}
