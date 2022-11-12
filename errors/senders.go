package errors

import (
	"errors"

	"github.com/alex-ilchukov/flow/values"
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
// provided type E from [Senders] set. It returns the channels as error-writing
// channels in resulting array werrs of type E and _the same_ channels as
// error-reading channels in resulting slice rerrs. In special case of type E
// being a type of zero-elements array, it returns nil in rerrs.
func Make[E Senders]() (werrs E, rerrs []<-chan error) {
	l := len(werrs)
	if l == 0 {
		return
	}

	rerrs = make([]<-chan error, l, l)
	for i := 1; i <= l; i++ {
		c := make(chan error)
		// The trick is to cajole Go compiler. If the code is rewritten
		// as werrs[i] = c, it whines something regarding violation of
		// boundaries.
		werrs[l-i] = c
		rerrs[l-i] = c
	}

	return
}

// Close takes slice of error-writing channels and closes its every element.
func Close(werrs []chan<- error) {
	for _, ch := range werrs {
		close(ch)
	}
}
