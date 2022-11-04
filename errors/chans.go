package errors

// No is type of zero-elements arrays of error-writing channels, effectively
// meaning that its users are simple enough to need no error handling.
type No [0]chan<- error

// One is type of one-element arrays of error-writing channels used by regular
// cases.
type One [1]chan<- error

// Chans is set of common types of error-writing channels collections. The
// zero-elements array is for those cases, which don't need error handling at
// all, and one-element array is for regular cases.
type Chans interface {
	No | One
}

// Make creates as many error channels as there are elements in array of the
// provided type E from [Chans] set. It returns the channels as error-writing
// channels in resulting array werrs of type E and _the same_ channels as
// error-reading channels in resulting slice rerrs. In special case of type E
// being a type of zero-elements array, it returns nil in rerrs.
func Make[E Chans]() (werrs E, rerrs []<-chan error) {
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

// Closes takes array of error-writing channels of the provided type E from set
// [Chans] and closes its every element if there is any. It is supposed to work
// in pair with [Make].
func Close[E Chans](errs E) {
	for i, l := 1, len(errs); i <= l; i++ {
		close(errs[l-i])
	}
}
