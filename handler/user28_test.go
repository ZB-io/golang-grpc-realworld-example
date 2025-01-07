package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
)





type mockUserStore struct {
	getByIDFunc func(id uint) (*model.User, error)
}


/*
ROOST_METHOD_HASH=CurrentUser_e3fa631d55
ROOST_METHOD_SIG_HASH=CurrentUser_29413339e9

FUNCTION_DEF=func (h *Handler) CurrentUser(ctx context.Context, req *pb.Empty) (*pb.UserResponse, error) 

 */
func (m *mockLogger) Error() *zerolog.Event {
	return &zerolog.Event{}
}

func (m *mockUserStore) GetByID(id uint) (*model.User, error) {
	return m.getByIDFunc(id)
}

func (m *mockLogger) Info() *zerolog.Event {
	return &zerolog.Event{}
}

func TestHandlerCurrentUser(t *testing.T) {

	tests := []struct {
		name          string
		setupContext  func() context.Context
		setupMocks    func(*mockUserStore, *mockLogger)
		expectedResp  *pb.UserResponse
		expectedError error
	}{
		{
			name: "Successful Current User Retrieval",
			setupContext: func() context.Context {
				ctx := context.Background()

				return ctx
			},
			setupMocks: func(us *mockUserStore, l *mockLogger) {
				us.getByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{
						ID:       1,
						Email:    "test@example.com",
						Username: "testuser",
						Bio:      "test bio",
						Image:    "test.jpg",
					}, nil
				}
			},
			expectedResp: &pb.UserResponse{
				User: &pb.User{
					Email:    "test@example.com",
					Username: "testuser",
					Bio:      "test bio",
					Image:    "test.jpg",
					Token:    "mock-token",
				},
			},
			expectedError: nil,
		},
		{
			name: "Unauthenticated Request",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMocks:    func(us *mockUserStore, l *mockLogger) {},
			expectedResp:  nil,
			expectedError: status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Valid Token but User Not Found",
			setupContext: func() context.Context {
				ctx := context.Background()

				return ctx
			},
			setupMocks: func(us *mockUserStore, l *mockLogger) {
				us.getByIDFunc = func(id uint) (*model.User, error) {
					return nil, errors.New("user not found")
				}
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.NotFound, "user not found"),
		},
		{
			name: "Token Generation Failure",
			setupContext: func() context.Context {
				ctx := context.Background()

				return ctx
			},
			setupMocks: func(us *mockUserStore, l *mockLogger) {
				us.getByIDFunc = func(id uint) (*model.User, error) {
					return &model.User{ID: 1}, nil
				}

			},
			expectedResp:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
		},
		{
			name: "Context Cancellation",
			setupContext: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			setupMocks:    func(us *mockUserStore, l *mockLogger) {},
			expectedResp:  nil,
			expectedError: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockUS := &mockUserStore{}
			mockLogger := &mockLogger{}

			h := &Handler{
				logger: &zerolog.Logger{},
				us:     &store.UserStore{},
			}

			tt.setupMocks(mockUS, mockLogger)

			ctx := tt.setupContext()

			resp, err := h.CurrentUser(ctx, &pb.Empty{})

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResp.User.Email, resp.User.Email)
				assert.Equal(t, tt.expectedResp.User.Username, resp.User.Username)
				assert.NotEmpty(t, resp.User.Token)
			}

			t.Logf("Completed test: %s", tt.name)
		})
	}
}

