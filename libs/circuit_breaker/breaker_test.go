package circuit_breaker

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCircuitBreaker(t *testing.T) {
	type TestCase struct {
		name               string
		expectedStatusCode int
		before, after      func(ctx context.Context, args HookInput) (*HookOutput, error)
		serverHandler      func(w http.ResponseWriter, r *http.Request)
	}

	cases := []TestCase{
		{
			name:               "on post request, break & return 500",
			expectedStatusCode: http.StatusInternalServerError,
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				panic("this shouldn't be fired")
			},
			before: func(ctx context.Context, args HookInput) (*HookOutput, error) {
				if args.OutgoingRequest.Method == "POST" {
					return &HookOutput{
						IncomingResponse: &http.Response{
							Status:     "500 Internal Server Error",
							StatusCode: http.StatusInternalServerError,
							Header:     make(http.Header),
							Body:       io.NopCloser(strings.NewReader("500")),
							Request:    args.OutgoingRequest,
						},
					}, nil
				}
				return nil, nil
			},
			after: func(ctx context.Context, args HookInput) (*HookOutput, error) {
				return nil, nil
			},
		},
		{
			name:               "do nothing and forward to original server",
			expectedStatusCode: http.StatusInternalServerError,
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{}"))
			},
			before: func(ctx context.Context, args HookInput) (*HookOutput, error) {
				return nil, nil
			},
			after: func(ctx context.Context, args HookInput) (*HookOutput, error) {
				return nil, nil
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ts := httptest.NewServer(&mockServerHandler{
				t:       t,
				handler: c.serverHandler,
			})
			defer ts.Close()
			breaker := New(nil, &HookMock{before: c.before, after: c.after})

			client := http.Client{
				Transport: breaker,
			}

			resp, err := client.Post(ts.URL, "application/json", io.NopCloser(strings.NewReader("hello")))
			require.NoError(t, err)

			require.Equal(t, c.expectedStatusCode, resp.StatusCode)

		})
	}
}

type mockServerHandler struct {
	t       *testing.T
	handler func(w http.ResponseWriter, r *http.Request)
}

func (msh *mockServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msh.handler(w, r)
}

type HookMock struct {
	before, after func(ctx context.Context, args HookInput) (*HookOutput, error)
}

func (hm *HookMock) Before(ctx context.Context, args HookInput) (*HookOutput, error) {
	return hm.before(ctx, args)
}

func (hm *HookMock) After(ctx context.Context, args HookInput) (*HookOutput, error) {
	return hm.after(ctx, args)
}
