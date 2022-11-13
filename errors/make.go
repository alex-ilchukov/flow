package errors

import (
	"context"

	"github.com/alex-ilchukov/flow/values/sender"
)

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
