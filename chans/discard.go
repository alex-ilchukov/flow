package chans

// Discard reads values of type V in the blocking way from the provided channel
// til the channel is closed. It does nothing with the values.
func Discard[V any](c <-chan V) {
	for range c {
	}
}
