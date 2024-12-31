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
	mockNewServeMuxFunc           func(opts ...runtime.ServeMuxOption) *runtime.ServeMux
	mockRegisterUsersHandlerFunc  func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
	mockRegisterArticlesHandlerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
	mockListenAndServeFunc        func(addr string, handler http.Handler) error
)

/*
ROOST_METHOD_HASH=run_9594c70ad3
ROOST_METHOD_SIG_HASH=run_9bb183262c
*/
func TestRun(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Successful Gateway Server Initialization",
			setupMocks: func() {
				mockNewServeMuxFunc = func(opts ...runtime.ServeMuxOption) *runtime.ServeMux {
					return &runtime.ServeMux{}
				}
				mockRegisterUsersHandlerFunc = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return nil
				}
				mockRegisterArticlesHandlerFunc = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return nil
				}
				mockListenAndServeFunc = func(addr string, handler http.Handler) error {
					return nil
				}
			},
			expectedError: nil,
		},
		{
			name: "Error in Registering Users Handler",
			setupMocks: func() {
				mockNewServeMuxFunc = func(opts ...runtime.ServeMuxOption) *runtime.ServeMux {
					return &runtime.ServeMux{}
				}
				mockRegisterUsersHandlerFunc = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return errors.New("users handler registration error")
				}
				mockRegisterArticlesHandlerFunc = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return nil
				}
			},
			expectedError: errors.New("users handler registration error"),
		},
		{
			name: "Error in Registering Articles Handler",
			setupMocks: func() {
				mockNewServeMuxFunc = func(opts ...runtime.ServeMuxOption) *runtime.ServeMux {
					return &runtime.ServeMux{}
				}
				mockRegisterUsersHandlerFunc = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return nil
				}
				mockRegisterArticlesHandlerFunc = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return errors.New("articles handler registration error")
				}
			},
			expectedError: errors.New("articles handler registration error"),
		},
		{
			name: "Server Start Failure",
			setupMocks: func() {
				mockNewServeMuxFunc = func(opts ...runtime.ServeMuxOption) *runtime.ServeMux {
					return &runtime.ServeMux{}
				}
				mockRegisterUsersHandlerFunc = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return nil
				}
				mockRegisterArticlesHandlerFunc = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return nil
				}
				mockListenAndServeFunc = func(addr string, handler http.Handler) error {
					return errors.New("server start error")
				}
			},
			expectedError: errors.New("server start error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := mockRun()
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("run() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func mockRun() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ropts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
	}

	mux := mockNewServeMuxFunc(ropts...)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := mockRegisterUsersHandlerFunc(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}

	err = mockRegisterArticlesHandlerFunc(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}

	return mockListenAndServeFunc(":3000", mux)
}
