package merge

import (
	"context"
	"sync"

	"github.com/alex-ilchukov/flow/values"
)

// New creates new command on merging values from all provided channels to a
// new channel within the provided context.
func New[V any](ctx context.Context, chs []<-chan V) *cmd[V] {
	result := make(chan V)
	return &cmd[V]{ctx: ctx, chs: chs, result: result}
}

type cmd[V any] struct {
	ctx    context.Context
	chs    []<-chan V
	result chan V
	wg     sync.WaitGroup
}

// Call creates new channel and launches non-blocking concurrent reading of the
// channels, redirecting any appeared value to the new channel with respect of
// cancellation within the provided context. It returns the created channel.
// The function takes care of closing of the created channel in distinct
// go-routine automatically, when all the original channels are closed.
func (c *cmd[V]) Call() <-chan V {
	c.wg.Add(len(c.chs))

	for _, ch := range c.chs {
		go c.listen(ch)
	}

	go c.wait()

	return c.result
}

func (c *cmd[V]) listen(ch <-chan V) {
	defer c.wg.Done()

	for {
		v, status := values.Receive(c.ctx, ch)
		if status != nil {
			return
		}

		status = values.Send(c.ctx, c.result, v)
		if status != nil {
			return
		}
	}
}

func (c *cmd[_]) wait() {
	c.wg.Wait()
	close(c.result)
}
