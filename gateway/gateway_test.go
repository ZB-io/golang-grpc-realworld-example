package main

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	gw "github.com/raahii/golang-grpc-realworld-example/proto"
)

/*
ROOST_METHOD_HASH=run_9594c70ad3
ROOST_METHOD_SIG_HASH=run_9bb183262c


 */
func (s *mockArticlesServer) GetArticle(ctx context.Context, req *gw.GetArticleRequest) (*gw.GetArticleResponse, error) {
	return &gw.GetArticleResponse{}, nil
}

func (s *mockArticlesServer) ListArticles(ctx context.Context, req *gw.ListArticlesRequest) (*gw.ListArticlesResponse, error) {
	return &gw.ListArticlesResponse{}, nil
}

func (s *mockUsersServer) Login(ctx context.Context, req *gw.LoginRequest) (*gw.LoginResponse, error) {
	return &gw.LoginResponse{}, nil
}

func (s *mockUsersServer) Register(ctx context.Context, req *gw.RegisterRequest) (*gw.RegisterResponse, error) {
	return &gw.RegisterResponse{}, nil
}

func TestRun(t *testing.T) {
	tests := []struct {
		name            string
		setupMock       func(*testing.T)
		cleanup         func()
		expectedError   bool
		errorContains   string
		blockPort       bool
		cancelContext   bool
		invalidEndpoint bool
	}{
		{
			name: "Successful Server Initialization",
			setupMock: func(t *testing.T) {

				lis, err := net.Listen("tcp", ":50051")
				if err != nil {
					t.Fatal(err)
				}
				s := grpc.NewServer()

				gw.RegisterUsersServer(s, &mockUsersServer{})
				gw.RegisterArticlesServer(s, &mockArticlesServer{})
				go s.Serve(lis)
				t.Cleanup(func() {
					s.Stop()
					lis.Close()
				})
			},
			cleanup:       func() {},
			expectedError: false,
		},
		{
			name: "Context Cancellation",
			setupMock: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				t.Cleanup(cancel)
				go func() {
					time.Sleep(50 * time.Millisecond)
					cancel()
				}()
			},
			cleanup:       func() {},
			expectedError: true,
			cancelContext: true,
			errorContains: "context canceled",
		},
		{
			name: "Port Already In Use",
			setupMock: func(t *testing.T) {
				listener, err := net.Listen("tcp", ":3000")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					listener.Close()
				})
			},
			cleanup:       func() {},
			expectedError: true,
			blockPort:     true,
			errorContains: "address already in use",
		},
		{
			name: "Invalid Endpoint Configuration",
			setupMock: func(t *testing.T) {
				echoEndpoint = "invalid:endpoint"
			},
			cleanup: func() {
				echoEndpoint = "localhost:50051"
			},
			expectedError:   true,
			invalidEndpoint: true,
			errorContains:   "invalid endpoint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock(t)
			defer tt.cleanup()

			errChan := make(chan error, 1)

			go func() {
				err := run()
				errChan <- err
			}()

			var err error
			select {
			case err = <-errChan:
			case <-time.After(2 * time.Second):
				if tt.expectedError {
					t.Error("Expected error but got none")
				}
				return
			}

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got nil")
					return
				}
				if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing %q, got %v", tt.errorContains, err)
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr
}

