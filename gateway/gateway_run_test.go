package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/proto"
)

type Controller struct {
	// T should only be called within a generated mock. It is not intended to
	// be used in user code and may be changed in future versions. T is the
	// TestReporter passed in when creating the Controller via NewController.
	// If the TestReporter does not implement a TestHelper it will be wrapped
	// with a nopTestHelper.
	T             TestHelper
	mu            sync.Mutex
	expectedCalls *callSet
	finished      bool
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func Testrun(t *testing.T) {

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
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			proto.UsersHandler.RegisterFunc = tt.mockUsersReg
			proto.ArticlesHandler.RegisterFunc = tt.mockArticlesReg
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
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}

			t.Log("Test", tt.name, "completed with error:", err)
		})
	}
}
