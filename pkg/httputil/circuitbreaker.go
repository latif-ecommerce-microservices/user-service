package httputil

import "github.com/sony/gobreaker/v2"

type CircuitBreakerer interface {
	Execute(fn func() (any, error)) (any, error)
}

type CircuitBreaker struct {
	gobreaker *gobreaker.CircuitBreaker[any]
}

type CircuitBreakerSettings gobreaker.Settings

func (cb *CircuitBreaker) Execute(fn func() (any, error)) (any, error) {
	return cb.gobreaker.Execute(fn)
}

type NoopCircuitBreaker struct{}

func (cb *NoopCircuitBreaker) Execute(fn func() (any, error)) (any, error) {
	return fn()
}
