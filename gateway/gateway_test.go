package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc"
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

