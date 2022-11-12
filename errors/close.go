package errors

// Close takes slice of error-writing channels and closes its every element.
func Close(werrs []chan<- error) {
	for _, ch := range werrs {
		close(ch)
	}
}
