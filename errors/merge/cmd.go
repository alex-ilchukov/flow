package merge

import (
	"context"
	"sync"

	"github.com/alex-ilchukov/flow/values"
)

type cmd struct {
	ctx    context.Context
	errs   []<-chan error
	wg     sync.WaitGroup
	result chan error
}

// New creates new command on merging errors from all error channels from the
// provided collection errs to a new error channel within the provided context.
func New(ctx context.Context, errs ...[]<-chan error) *cmd {
	flattened := make([]<-chan error, 0)
	for _, e := range errs {
		flattened = append(flattened, e...)
	}

	result := make(chan error)
	return &cmd{ctx: ctx, errs: flattened, result: result}
}

// Call creates new error channel and launches non-blocking concurrent reading
// of the channels, redirecting any appeared error to the new error channel
// with respect of cancellation within the provided context. In the end it
// returns the error channel. It takes care of closing of the error channel in
// distinct go-routine, so a user should not be bothered by it.
func (c *cmd) Call() <-chan error {
	c.wg.Add(len(c.errs))

	for _, ch := range c.errs {
		go c.listen(ch)
	}

	go c.wait()

	return c.result
}

func (c *cmd) listen(ch <-chan error) {
	defer c.wg.Done()

	for {
		err, status := values.Receive(c.ctx, ch)
		if status != nil {
			return
		}

		status = values.Send(c.ctx, c.result, err)
		if status != nil {
			return
		}
	}
}

func (c *cmd) wait() {
	c.wg.Wait()
	close(c.result)
}
