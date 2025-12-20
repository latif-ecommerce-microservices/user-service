package container

import (
	"net/http"
	"os"
	"time"
)

type Option func(*Container)

// WithTimeout is an option to provide custom timeout when shutdown the server.
// The default timeout is 30 seconds if you don't provide a timeout number.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Container) {
		c.timeout = timeout
	}
}

// WithSignal allows you to define custom OS signal without relying on the default ones.
func WithSignal(sig <-chan os.Signal) Option {
	return func(c *Container) {
		c.signal = sig
	}
}

// WithHTTPServer allows you to add HTTP server in the App container.
// By default, the container doesn't provide HTTP server, instead the client should add it by themselves.
// That add flexibility to add the component or not since the application not always have a HTTP server.
func WithHTTPServer(srv []*http.Server) Option {
	return func(c *Container) {
		c.srv = srv
	}
}
