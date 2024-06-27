package circuit_breaker

import (
	"context"
	"net/http"
)

type HookOutput struct {
	IncomingResponse *http.Response
}

type HookInput struct {
	OutgoingRequest  *http.Request
	IncomingResponse *http.Response
}

// Hook provides a standardized interface for tapping into the request lifecycle.
type Hook interface {
	Before(ctx context.Context, args HookInput) (*HookOutput, error)
	After(ctx context.Context, args HookInput) (*HookOutput, error)
}

// CircuitBreaker provides a custom transport that allows programmatic simulation of api failures
type CircuitBreaker struct {
	rt   http.RoundTripper
	hook Hook
}

func New(rt http.RoundTripper, hook Hook) *CircuitBreaker {
	return &CircuitBreaker{
		rt:   rt,
		hook: hook,
	}
}

func (b *CircuitBreaker) RoundTrip(req *http.Request) (*http.Response, error) {
	if b.rt == nil {
		b.rt = http.DefaultTransport
	}

	hr, err := b.hook.Before(req.Context(), HookInput{
		OutgoingRequest: req,
	})
	if err != nil {
		return nil, err
	}
	if hr != nil {
		return hr.IncomingResponse, nil
	}

	resp, err := b.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	hr, err = b.hook.After(req.Context(), HookInput{
		OutgoingRequest:  req,
		IncomingResponse: resp,
	})
	if err != nil {
		return nil, err
	}
	if hr != nil {
		return hr.IncomingResponse, nil
	}

	return resp, err
}
