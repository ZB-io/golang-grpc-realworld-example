
// ********RoostGPT********
/*

roost_feedback [1/2/2025, 11:46:18 AM]:Need to Improve some test of this

roost_feedback [1/3/2025, 3:48:48 AM]:Need to Improve some test of this
*/

// ********RoostGPT********

package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type Controller struct {
	mu            sync.Mutex
	t             gomock.TestReporter
	expectedCalls *gomock.CallSet
	finished      bool
}

func TestRun(t *testing.T) {
	originalRegisterUsersHandlerFromEndpoint := proto.RegisterUsersHandlerFromEndpoint
	originalRegisterArticlesHandlerFromEndpoint := proto.RegisterArticlesHandlerFromEndpoint
	originalListenAndServe := http.ListenAndServe

	defer func() {
		proto.RegisterUsersHandlerFromEndpoint = originalRegisterUsersHandlerFromEndpoint
		proto.RegisterArticlesHandlerFromEndpoint = originalRegisterArticlesHandlerFromEndpoint
		http.ListenAndServe = originalListenAndServe
	}()

	tests := []struct {
		name            string
		mockUsersReg    func() error
		mockArticlesReg func() error
		mockServer      func() error
		endpoint        string
		expectedError   error
	}{
		{
			name: "Successful Gateway Initialization and Server Start",
			mockUsersReg: func() error {
				return nil
			},
			mockArticlesReg: func() error {
				return nil
			},
			mockServer: func() error {
				return nil
			},
			endpoint:      "localhost:50051",
			expectedError: nil,
		},
		{
			name: "Failure During UsersHandler Registration",
			mockUsersReg: func() error {
				return errors.New("failed to register Users handler")
			},
			mockArticlesReg: func() error {
				return nil
			},
			mockServer: func() error {
				return nil
			},
			endpoint:      "localhost:50051",
			expectedError: errors.New("failed to register Users handler"),
		},
		{
			name: "Failure During ArticlesHandler Registration",
			mockUsersReg: func() error {
				return nil
			},
			mockArticlesReg: func() error {
				return errors.New("failed to register Articles handler")
			},
			mockServer: func() error {
				return nil
			},
			endpoint:      "localhost:50051",
			expectedError: errors.New("failed to register Articles handler"),
		},
		{
			name: "HTTP Server Startup Error",
			mockUsersReg: func() error {
				return nil
			},
			mockArticlesReg: func() error {
				return nil
			},
			mockServer: func() error {
				return errors.New("failed to start HTTP server")
			},
			endpoint:      "localhost:50051",
			expectedError: errors.New("failed to start HTTP server"),
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			proto.RegisterUsersHandlerFromEndpoint = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return tt.mockUsersReg()
			}
			proto.RegisterArticlesHandlerFromEndpoint = func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
				return tt.mockArticlesReg()
			}
			http.ListenAndServe = func(addr string, handler http.Handler) error {
				return tt.mockServer()
			}

			oldArgs := flag.Args()
			defer flag.CommandLine.Parse(oldArgs)
			flag.CommandLine.Parse([]string{"-endpoint", tt.endpoint})

			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(nil)
			}()

			err := run()

			if tt.expectedError != nil {
				assert.NotNil(t, err, "Expected error but got none")
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Contains(t, buf.String(), tt.expectedError.Error(), "Expected log message containing error")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Contains(t, buf.String(), "Server started", "Expected log message indicating server start")
			}

			t.Log("Test", tt.name, "completed with error:", err)
		})
	}
}
