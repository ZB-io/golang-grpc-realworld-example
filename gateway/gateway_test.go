package main

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	mockRegisterUsersHandlerFromEndpoint    func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
	mockRegisterArticlesHandlerFromEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
	mockListenAndServe                      func(addr string, handler http.Handler) error
	mockNewServeMux                         func(opts ...runtime.ServeMuxOption) *runtime.ServeMux
)

/*
ROOST_METHOD_HASH=run_9594c70ad3
ROOST_METHOD_SIG_HASH=run_9bb183262c
*/
func TestRun(t *testing.T) {
	originalRegisterUsersHandler := RegisterUsersHandlerFromEndpoint
	originalRegisterArticlesHandler := RegisterArticlesHandlerFromEndpoint
	originalListenAndServe := http.ListenAndServe
	originalNewServeMux := runtime.NewServeMux

	defer func() {
		RegisterUsersHandlerFromEndpoint = originalRegisterUsersHandler
		RegisterArticlesHandlerFromEndpoint = originalRegisterArticlesHandler
		http.ListenAndServe = originalListenAndServe
		runtime.NewServeMux = originalNewServeMux
	}()

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful Gateway Server Initialization",
			setupMocks: func() {
				mockRegisterUsersHandlerFromEndpoint = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return nil
				}
				mockRegisterArticlesHandlerFromEndpoint = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return nil
				}
				mockListenAndServe = func(addr string, handler http.Handler) error {
					return nil
				}
				mockNewServeMux = func(opts ...runtime.ServeMuxOption) *runtime.ServeMux {
					return &runtime.ServeMux{}
				}
			},
			expectedError: nil,
		},
		{
			name: "Error in Registering Users Handler",
			setupMocks: func() {
				mockRegisterUsersHandlerFromEndpoint = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return errors.New("users handler error")
				}
				mockRegisterArticlesHandlerFromEndpoint = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return nil
				}
				mockListenAndServe = func(addr string, handler http.Handler) error {
					return nil
				}
				mockNewServeMux = func(opts ...runtime.ServeMuxOption) *runtime.ServeMux {
					return &runtime.ServeMux{}
				}
			},
			expectedError: errors.New("users handler error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			RegisterUsersHandlerFromEndpoint = mockRegisterUsersHandlerFromEndpoint
			RegisterArticlesHandlerFromEndpoint = mockRegisterArticlesHandlerFromEndpoint
			http.ListenAndServe = mockListenAndServe
			runtime.NewServeMux = mockNewServeMux

			err := run()

			if (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError == nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("run() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
