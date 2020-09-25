package closer

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"go.uber.org/multierr"
)

// Closer is graceful close all funcs.
type Closer struct {
	sync.Mutex
	single sync.Once
	funcs  []func() error
	done   chan struct{}
}

// New ...
func New(signals ...os.Signal) *Closer {
	c := &Closer{
		done: make(chan struct{}),
	}
	if len(signals) == 0 {
		return c
	}

	go func() {
		signalsCh := make(chan os.Signal, 1)
		signal.Notify(signalsCh, signals...)
		<-signalsCh
		signal.Stop(signalsCh)
		c.Close()
	}()

	return c
}

// Close ...
func (c *Closer) Close() {
	c.single.Do(func() {
		defer close(c.done)
	})

	c.Lock()
	funcs := c.funcs
	c.funcs = nil
	c.Unlock()

	errs := make(chan error, len(funcs))
	for _, f := range funcs {
		go func(f func() error) {
			errs <- f()
		}(f)
	}

	var resultErr error
	for err := range errs {
		if err != nil {
			resultErr = multierr.Append(resultErr, err)
		}
	}
	if resultErr != nil {
		_, _ = fmt.Fprintln(os.Stdout, resultErr)
	}

	return
}

// Bind add close func.
func (c *Closer) Bind(f func() error) {
	c.Lock()
	c.funcs = append(c.funcs, f)
	c.Unlock()
}
