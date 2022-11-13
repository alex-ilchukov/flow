package errors_test

import (
	"context"
	"errors"
	"testing"
	"time"

	flowerrors "github.com/alex-ilchukov/flow/errors"
)

func TestMakeNoErrors(t *testing.T) {
	ctx := context.Background()
	_, rerrs, werrs := flowerrors.Make[flowerrors.No](ctx)
	switch {
	case len(rerrs) > 0:
		t.Errorf("Got something in rerrs")

	case len(werrs) > 0:
		t.Errorf("Got something in werrs")
	}
}

func TestMakeOneError(t *testing.T) {
	ctx := context.Background()
	s, rerrs, werrs := flowerrors.Make[flowerrors.One](ctx)
	switch {
	case len(rerrs) != 1:
		t.Errorf("Invalid length of rerrs")

	case len(rerrs) != len(werrs):
		t.Errorf("Lengths of rerrs and werrs aren't the same")
	}

	deadline := time.Now().Add(time.Millisecond)
	ctx, _ = context.WithDeadline(ctx, deadline)
	err := errors.New("test")
	go func() { s[0].Send(err) }()

	select {
	case e := <-rerrs[0]:
		if e != err {
			t.Errorf("Invalid error %v", e)
		}

	case <-ctx.Done():
		t.Errorf("Got no error")
	}
}
