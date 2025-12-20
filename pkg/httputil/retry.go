package httputil

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v5"
)

type (
	RetryOperation = backoff.Operation[any]
	RetryOption    = backoff.RetryOption
	Notify         = backoff.Notify

	Retrier interface {
		Retry(context.Context, RetryOperation, ...RetryOption) (any, error)
	}

	NoopRetry               struct{}
	ExponentialBackoffRetry struct{}
	ConstantRetry           struct {
		Interval time.Duration
	}
)

var (
	RetryAfter         = backoff.RetryAfter
	PermanentError     = backoff.Permanent
	WithMaxElapsedTime = backoff.WithMaxElapsedTime
	WithMaxTries       = backoff.WithMaxTries
	WithNotify         = backoff.WithNotify
)

func NewNoopRetry() *NoopRetry {
	return &NoopRetry{}
}

func (n *NoopRetry) Retry(ctx context.Context, operation RetryOperation, opts ...RetryOption) (any, error) {
	return operation()
}

func NewExponentialBackOffRetry() *ExponentialBackoffRetry {
	return &ExponentialBackoffRetry{}
}

func (e *ExponentialBackoffRetry) Retry(ctx context.Context, operation RetryOperation, opts ...RetryOption) (any, error) {
	ropts := append(opts, backoff.WithBackOff(backoff.NewExponentialBackOff()))

	return backoff.Retry(ctx, operation, ropts...)
}

func NewConstantRetry() *ConstantRetry {
	return &ConstantRetry{}
}

func (c *ConstantRetry) Retry(ctx context.Context, operation RetryOperation, opts ...RetryOption) (any, error) {
	ropts := append(opts, backoff.WithBackOff(backoff.NewConstantBackOff(c.Interval)))

	return backoff.Retry(ctx, operation, ropts...)
}
