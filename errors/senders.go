package errors

import (
	"context"
	"errors"

	"github.com/alex-ilchukov/flow/values"
	"github.com/alex-ilchukov/flow/values/sender"
)

// No is type of zero-elements arrays of error senders, effectively meaning
// that its users are simple enough to need no error reporting.
type No [0]values.Sender[error]

// One is type of one-element arrays of error senders used by regular cases.
type One [1]values.Sender[error]

// Send takes error value and tries to send it, delegating the process to the
// only element of a. It returns nil, if no error appears, or the error
// otherwise. Additionally, it returns error value of the element is nil.
func (a *One) Send(err error) error {
	if a == nil || a[0] == nil {
		return errors.New("nil element")
	}

	return a[0].Send(err)
}

// Senders is set of common types of error-senders collections. The
// zero-elements array is for those cases, which don't need error handling at
// all, and one-element array is for regular cases.
type Senders interface {
	No | One
}

// Make creates as many error channels as there are elements in array of the
// provided type E from [Senders] set and setups the elements as error senders
// with the channels and the provided context. It returns the array and the
// channels in form of slices of error-reading and error-writing channels. In
// special case of type E being a type of zero-elements array, it returns nil
// in rerrs and werrs.
func Make[E Senders](ctx context.Context) (
	s E,
	rerrs []<-chan error,
	werrs []chan<- error,
) {

	l := len(s)
	if l == 0 {
		return
	}

	rerrs = make([]<-chan error, l, l)
	werrs = make([]chan<- error, l, l)
	for i := 1; i <= l; i++ {
		c := make(chan error)
		// The trick is to cajole Go compiler. If the code is rewritten
		// as s[i] = sender.New(c) with i going from 0 to 0, the
		// compiler whines on something regarding violation of
		// boundaries.
		s[l-i] = sender.New(ctx, c)
		rerrs[l-i] = c
		werrs[l-i] = c
	}

	return
}
