package errors_test

import (
	"context"
	"errors"
	"testing"
	"time"

	flowerrors "github.com/alex-ilchukov/flow/errors"
)

func TestMakeNoErrors(t *testing.T) {
	_, rerrs := flowerrors.Make[flowerrors.No]()
	if len(rerrs) > 0 {
		t.Errorf("Got something in rerrs")
	}
}

func TestMakeOneError(t *testing.T) {
	werrs, rerrs := flowerrors.Make[flowerrors.One]()
	if len(rerrs) != 1 {
		t.Errorf("Invalid length of rerrs")
	}

	ctx := context.Background()
	deadline := time.Now().Add(time.Millisecond)
	ctx, _ = context.WithDeadline(ctx, deadline)
	err := errors.New("test")
	go func() { werrs[0] <- err }()

	select {
	case e := <-rerrs[0]:
		if e != err {
			t.Errorf("Invalid error %v", e)
		}

	case <-ctx.Done():
		t.Errorf("Got no error")
	}
}
