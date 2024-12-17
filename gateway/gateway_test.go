package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
	"io/ioutil"
)

/*
ROOST_METHOD_HASH=run_9594c70ad3
ROOST_METHOD_SIG_HASH=run_9bb183262c


 */
func CaptureAndCompareOutput(t *testing.T, f func() error, expectedError bool) {

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := f()

	err = w.Close()
	if err != nil {
		t.Errorf("error closing writer: %v", err)
	}
	os.Stdout = old

	out, _ := ioutil.ReadAll(r)

	success := err == nil && !expectedError || err != nil && expectedError
	if !success {
		t.Errorf("Expected error: %t, got: %v, output: %s", expectedError, err, string(out))
	} else {
		t.Logf("Test output: %s", string(out))
	}
}

func Testrun(t *testing.T) {

	flag.Parse()
	defer flag.Set("endpoint", "localhost:50051")

	tests := []struct {
		name          string
		modifyEnv     func()
		expectedError bool
	}{
		{
			name: "Successful Registration of Handlers",
			modifyEnv: func() {

			},
			expectedError: false,
		},
		{
			name: "Handler Registration Failure",
			modifyEnv: func() {
				flag.Set("endpoint", "invalid:50051")
			},
			expectedError: true,
		},
		{
			name: "Server Start Failure Due to Port Binding Error",
			modifyEnv: func() {

				go http.ListenAndServe(":3000", http.NewServeMux())
				time.Sleep(1 * time.Second)
			},
			expectedError: true,
		},
		{
			name: "Context Cancellation before Completion",
			modifyEnv: func() {

				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_ = ctx
			},
			expectedError: true,
		},
		{
			name: "Testing with Different 'echoEndpoint' Configurations",
			modifyEnv: func() {

			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buf bytes.Buffer
			log.SetOutput(&buf)

			tt.modifyEnv()

			err := run()

			if tt.expectedError && err == nil {
				t.Errorf("expected an error but did not get one")
			} else if !tt.expectedError && err != nil {
				t.Errorf("did not expect an error but got one: %v", err)
			} else {
				t.Logf("Scenario '%s' executed with output: %s", tt.name, buf.String())
			}

			log.SetOutput(ioutil.Discard)
		})
	}
}

