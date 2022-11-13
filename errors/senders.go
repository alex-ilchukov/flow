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
