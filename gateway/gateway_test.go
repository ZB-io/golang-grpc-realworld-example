package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gw "github.com/raahii/golang-grpc-realworld-example/proto"
)

var (
	mockRegisterUsersHandler    func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
	mockRegisterArticlesHandler func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
	mockListenAndServe          func(addr string, handler http.Handler) error
)

// Mock functions to replace the actual functions during testing
func mockRegisterUsersHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return mockRegisterUsersHandler(ctx, mux, endpoint, opts)
}

func mockRegisterArticlesHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return mockRegisterArticlesHandler(ctx, mux, endpoint, opts)
}

func mockHTTPListenAndServe(addr string, handler http.Handler) error {
	return mockListenAndServe(addr, handler)
}

/*
ROOST_METHOD_HASH=run_9594c70ad3
ROOST_METHOD_SIG_HASH=run_9bb183262c
*/
func TestRun(t *testing.T) {
	originalEchoEndpoint := *echoEndpoint
	defer func() {
		*echoEndpoint = originalEchoEndpoint
	}()

	// Save original functions and restore them after the test
	originalRegisterUsersHandler := gw.RegisterUsersHandlerFromEndpoint
	originalRegisterArticlesHandler := gw.RegisterArticlesHandlerFromEndpoint
	originalListenAndServe := http.ListenAndServe
	defer func() {
		gw.RegisterUsersHandlerFromEndpoint = originalRegisterUsersHandler
		gw.RegisterArticlesHandlerFromEndpoint = originalRegisterArticlesHandler
		http.ListenAndServe = originalListenAndServe
	}()

	tests := []struct {
		name                   string
		echoEndpoint           string
		mockUsersHandler       func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
		mockArticlesHandler    func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
		mockListenAndServe     func(addr string, handler http.Handler) error
		expectedError          error
		expectedUsersCall      bool
		expectedArticlesCall   bool
		expectedListenAndServe bool
	}{
		{
			name:         "Success",
			echoEndpoint: "localhost:50051",
			mockUsersHandler: func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return nil
			},
			mockArticlesHandler: func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return nil
			},
			mockListenAndServe: func(addr string, handler http.Handler) error {
				return nil
			},
			expectedError:          nil,
			expectedUsersCall:      true,
			expectedArticlesCall:   true,
			expectedListenAndServe: true,
		},
		{
			name:         "UsersHandlerError",
			echoEndpoint: "localhost:50051",
			mockUsersHandler: func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return errors.New("users handler error")
			},
			expectedError:     errors.New("users handler error"),
			expectedUsersCall: true,
		},
		{
			name:         "ArticlesHandlerError",
			echoEndpoint: "localhost:50051",
			mockUsersHandler: func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return nil
			},
			mockArticlesHandler: func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return errors.New("articles handler error")
			},
			expectedError:        errors.New("articles handler error"),
			expectedUsersCall:    true,
			expectedArticlesCall: true,
		},
		{
			name:         "ListenAndServeError",
			echoEndpoint: "localhost:50051",
			mockUsersHandler: func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return nil
			},
			mockArticlesHandler: func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return nil
			},
			mockListenAndServe: func(addr string, handler http.Handler) error {
				return errors.New("listen and serve error")
			},
			expectedError:          errors.New("listen and serve error"),
			expectedUsersCall:      true,
			expectedArticlesCall:   true,
			expectedListenAndServe: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*echoEndpoint = tt.echoEndpoint

			mockRegisterUsersHandler = tt.mockUsersHandler
			mockRegisterArticlesHandler = tt.mockArticlesHandler
			mockListenAndServe = tt.mockListenAndServe

			// Replace the actual functions with mock functions
			gw.RegisterUsersHandlerFromEndpoint = mockRegisterUsersHandlerFromEndpoint
			gw.RegisterArticlesHandlerFromEndpoint = mockRegisterArticlesHandlerFromEndpoint
			http.ListenAndServe = mockHTTPListenAndServe

			err := run()

			if (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError == nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("run() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
