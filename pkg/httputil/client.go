package httputil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cenkalti/backoff/v5"
	"github.com/sony/gobreaker/v2"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client interface {
	Get(ctx context.Context, url string, headers http.Header) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error)
	Put(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error)
	Patch(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error)
	Delete(ctx context.Context, url string, headers http.Header) (*http.Response, error)
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

type Option func(*HTTPClient)

type ErrHandle func(*http.Request, any) (any, error)

type HTTPClient struct {
	client     Doer
	cb         CircuitBreakerer
	retry      Retrier
	retryOpts  []RetryOption
	ErrHandler ErrHandle
}

func defaultErrHandle(req *http.Request, res any) (any, error) {
	httpRes := res.(*http.Response)

	if httpRes.StatusCode == 400 {
		return nil, backoff.Permanent(errors.New("http client: bad request"))
	}

	if httpRes.StatusCode >= 500 {
		return nil, fmt.Errorf("http client: %v", httpRes.Status)
	}

	if httpRes.StatusCode == 429 {
		seconds, err := strconv.ParseInt(httpRes.Header.Get("Retry-After"), 10, 64)
		if err == nil {
			return nil, backoff.RetryAfter(int(seconds))
		}
	}

	return res, nil
}

func NewHTTPClient(opts ...Option) *HTTPClient {
	httpClient := &HTTPClient{
		cb:         &NoopCircuitBreaker{},
		retry:      &NoopRetry{},
		retryOpts:  []RetryOption{},
		ErrHandler: defaultErrHandle,
	}

	for _, opt := range opts {
		opt(httpClient)
	}

	if httpClient.client == nil {
		httpClient.client = &http.Client{}
	}

	return httpClient
}

func (h *HTTPClient) Get(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GET request creation failed: %v", err)
	}

	request.Header = headers

	return h.Do(ctx, request)
}

func (h *HTTPClient) Post(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("POST request creation failed: %v", err)
	}

	request.Header = headers

	return h.Do(ctx, request)
}

func (h *HTTPClient) Put(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, fmt.Errorf("PUT request creation failed: %v", err)
	}

	request.Header = headers

	return h.Do(ctx, request)
}

func (h *HTTPClient) Patch(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, fmt.Errorf("PATCH request creation failed: %v", err)
	}

	request.Header = headers

	return h.Do(ctx, request)
}

func (h *HTTPClient) Delete(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("GET request creation failed: %v", err)
	}

	request.Header = headers

	return h.Do(ctx, request)
}

func (h *HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	operation := func() (any, error) {
		var err error
		resp, err := h.client.Do(req)
		if err != nil {
			return nil, err
		}

		return h.ErrHandler(req, resp)
	}

	res, err := h.cb.Execute(func() (any, error) {
		result, err := h.retry.Retry(ctx, operation, h.retryOpts...)
		if err != nil {
			return nil, err
		}

		return result, nil
	})
	if err != nil {
		return nil, err
	}

	return res.(*http.Response), nil
}

func WithCircuitBreaker(st CircuitBreakerSettings) Option {
	return func(h *HTTPClient) {
		settings := gobreaker.Settings{
			Name:          st.Name,
			MaxRequests:   st.MaxRequests,
			Interval:      st.Interval,
			Timeout:       st.Timeout,
			ReadyToTrip:   st.ReadyToTrip,
			OnStateChange: st.OnStateChange,
			IsSuccessful:  st.IsSuccessful,
		}

		h.cb = &CircuitBreaker{
			gobreaker: gobreaker.NewCircuitBreaker[any](settings),
		}
	}
}

func WithRetry(strategy Retrier) Option {
	return func(h *HTTPClient) {
		h.retry = strategy
	}
}

func WithRetryOptions(opts ...RetryOption) Option {
	return func(h *HTTPClient) {
		h.retryOpts = opts
	}
}

func WithCustomErrHandling(handle ErrHandle) Option {
	return func(h *HTTPClient) {
		h.ErrHandler = handle
	}
}
