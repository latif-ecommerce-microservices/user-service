package container

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App interface {
	BeforeStart(context.Context) error
	AfterStart(context.Context) error
}

type Container struct {
	App

	srv     []*http.Server
	signal  <-chan os.Signal
	timeout time.Duration
}

var (
	defaultSignal = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}

	defaultTimeoutDuration = 30 * time.Second
)

func newDefaultSigChannel() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, defaultSignal...)
	return sig
}
func New(app App, opts ...Option) *Container {
	container := &Container{
		App:     app,
		timeout: defaultTimeoutDuration,
		signal:  newDefaultSigChannel(),
	}

	for _, opt := range opts {
		opt(container)
	}

	return container
}

func (c *Container) Start(ctx context.Context) (err error) {

	if err = c.App.BeforeStart(ctx); err != nil {
		return
	}

	defer func() {
		stopErr := c.App.AfterStart(ctx)
		if stopErr != nil {
			if err == nil {
				err = stopErr
			} else {
				err = fmt.Errorf("%v; additionally failed to stop app: %v", err, stopErr)
			}
		}
	}()

	select {

	case err, ok := <-c.runServer():
		if ok && err != nil {
			return err
		}

	case <-ctx.Done():
		return c.shutdownServer(ctx)

	case <-c.signal:
		return c.shutdownServer(ctx)
	}

	return nil
}

func (c *Container) runServer() chan error {
	errCh := make(chan error, len(c.srv))

	if len(c.srv) == 0 {
		close(errCh)
		return errCh
	}

	for _, srv := range c.srv {
		go func(s *http.Server) {
			if err := s.ListenAndServe(); err != nil &&
				!errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}(srv)
	}

	return errCh
}

func (c *Container) shutdownServer(ctx context.Context) error {
	if len(c.srv) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	for _, srv := range c.srv {
		if err := srv.Shutdown(ctx); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
