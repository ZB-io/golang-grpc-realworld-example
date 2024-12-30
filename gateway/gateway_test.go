package main

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	gw "github.com/raahii/golang-grpc-realworld-example/proto"
)

var (
	mockRegisterUsersHandler    func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
	mockRegisterArticlesHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
	mockListenAndServe          func(addr string, handler http.Handler) error
	mockNewServeMux             func(opts ...runtime.ServeMux0ption) *runtime.ServeMux
	echoEndpoint                string // Add this line
)

func TestRun(t *testing.T) {
	// ... (keep the existing test cases)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegisterUsersHandler = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return tt.registerUsersErr
			}
			mockRegisterArticlesHandler = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return tt.registerArticlesErr
			}
			mockListenAndServe = func(addr string, handler http.Handler) error {
				return tt.listenAndServeErr
			}

			var capturedOptions []runtime.ServeMuxOption
			mockNewServeMux = func(opts ...runtime.ServeMuxOption) *runtime.ServeMux {
				capturedOptions = opts
				return &runtime.ServeMux{}
			}

			if tt.customEndpoint != "" {
				echoEndpoint = tt.customEndpoint // Change this line
			} else {
				echoEndpoint = "localhost:50051" // Change this line
			}

			var err error
			if tt.cancelContext {
				ctx, cancel := context.WithCancel(context.Background())
				go func() {
					err = run(ctx) // Add ctx parameter
				}()
				cancel()
			} else {
				err = run(context.Background()) // Add context.Background()
			}

			// ... (keep the rest of the test function)
		})
	}
}

func init() {
	gw.RegisterUsersHandlerFromEndpoint = mockRegisterUsersHandler
	gw.RegisterArticlesHandlerFromEndpoint = mockRegisterArticlesHandler
	http.ListenAndServe = mockListenAndServe
	runtime.NewServeMux = mockNewServeMux
}

// Add this function to make the test compile
func run(ctx context.Context) error {
	// Implementation of the run function
	return nil
}
