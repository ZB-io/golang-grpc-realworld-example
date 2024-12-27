package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	gw "github.com/raahii/golang-grpc-realworld-example/proto"
)

/*
ROOST_METHOD_HASH=run_9594c70ad3
ROOST_METHOD_SIG_HASH=run_9bb183262c
*/
func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		configureFunc func() string
		expectedErr   error
		serverFunc    func(handler http.Handler) *httptest.Server
		contextCancel bool
	}{
		{
			name: "Scenario 1: Successful Registration of Handlers",
			configureFunc: func() string {
				listener := bufconn.Listen(1024 * 1024)
				server := grpc.NewServer()
				go func() {
					gw.RegisterUsersHandlerFunc(nil, nil, nil)
					gw.RegisterArticlesHandlerFunc(nil, nil, nil)
					if err := server.Serve(listener); err != nil {
						t.Fatalf("Server failed: %v", err)
					}
				}()
				return "bufnet://"
			},
			expectedErr: nil,
		},
		{
			name: "Scenario 2: Handler Registration Failure",
			configureFunc: func() string {
				return "invalid:50051"
			},
			expectedErr: errors.New("rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing dial tcp: address invalid: missing port in address\""),
		},
		{
			name: "Scenario 3: Server Start Failure Due to Port Binding Error",
			configureFunc: func() string {
				return "localhost:50051"
			},
			serverFunc: func(handler http.Handler) *httptest.Server {
				server := httptest.NewUnstartedServer(handler)
				server.Listener.Close()
				server.Config.Addr = ":3000"
				server.Start()
				return server
			},
			expectedErr: errors.New("listen tcp :3000: bind: address already in use"),
		},
		{
			name: "Scenario 4: Context Cancellation before Completion",
			configureFunc: func() string {
				return "localhost:50051"
			},
			contextCancel: true,
			expectedErr:   errors.New("context canceled"),
		},
		{
			name: "Scenario 5: Testing with Different `echoEndpoint` Configurations",
			configureFunc: func() string {
				return "localhost:50052"
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.serverFunc != nil {
				server := tt.serverFunc(nil)
				if server != nil {
					defer server.Close()
				}
			}

			*echoEndpoint = tt.configureFunc()

			ctx := context.Background()
			if tt.contextCancel {
				cancelCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
				ctx = cancelCtx
				defer cancel()
			}

			var capturedErr error
			go func(c *http.Client) {
				err := run()
				capturedErr = err
			}(http.DefaultClient)

			time.Sleep(200 * time.Millisecond)

			if (tt.expectedErr == nil && capturedErr != nil) || (tt.expectedErr != nil && capturedErr == nil) ||
				(tt.expectedErr != nil && capturedErr != nil && tt.expectedErr.Error() != capturedErr.Error()) {
				t.Errorf("Test %s failed: expected error '%v', got error '%v'", tt.name, tt.expectedErr, capturedErr)
			}
		})
	}
}
