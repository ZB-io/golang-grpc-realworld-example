package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type MockArticleStore struct {
	mock.Mock
}

/*
ROOST_METHOD_HASH=GetTags_42221e4328
ROOST_METHOD_SIG_HASH=GetTags_52f72598a3
*/
func (m *MockArticleStore) GetTags() ([]model.Tag, error) {
	args := m.Called()
	return args.Get(0).([]model.Tag), args.Error(1)
}

func TestGetTags(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockArticleStore)
		setupContext   func() (context.Context, context.CancelFunc)
		expectedTags   []string
		expectedError  error
		expectedStatus codes.Code
	}{
		{
			name: "Successful Retrieval of Tags",
			setupMock: func(mas *MockArticleStore) {
				mas.On("GetTags").Return([]model.Tag{{Name: "tag1"}, {Name: "tag2"}}, nil)
			},
			setupContext: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			expectedTags:   []string{"tag1", "tag2"},
			expectedError:  nil,
			expectedStatus: codes.OK,
		},
		{
			name: "Empty Tags List",
			setupMock: func(mas *MockArticleStore) {
				mas.On("GetTags").Return([]model.Tag{}, nil)
			},
			setupContext: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			expectedTags:   []string{},
			expectedError:  nil,
			expectedStatus: codes.OK,
		},
		{
			name: "Error from Article Store",
			setupMock: func(mas *MockArticleStore) {
				mas.On("GetTags").Return([]model.Tag{}, errors.New("store error"))
			},
			setupContext: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			expectedTags:   nil,
			expectedError:  status.Error(codes.Aborted, "internal server error"),
			expectedStatus: codes.Aborted,
		},
		{
			name: "Context Cancellation",
			setupMock: func(mas *MockArticleStore) {
				mas.On("GetTags").After(100*time.Millisecond).Return([]model.Tag{}, nil)
			},
			setupContext: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				go func() {
					time.Sleep(50 * time.Millisecond)
					cancel()
				}()
				return ctx, cancel
			},
			expectedTags:   nil,
			expectedError:  context.Canceled,
			expectedStatus: codes.Canceled,
		},
		{
			name: "Large Number of Tags",
			setupMock: func(mas *MockArticleStore) {
				largeTags := make([]model.Tag, 10000)
				for i := 0; i < 10000; i++ {
					largeTags[i] = model.Tag{Name: "tag" + string(rune(i))}
				}
				mas.On("GetTags").Return(largeTags, nil)
			},
			setupContext: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			expectedTags:   nil,
			expectedError:  nil,
			expectedStatus: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockArticleStore)
			tt.setupMock(mockStore)

			logger := zerolog.New(zerolog.ConsoleWriter{Out: zerolog.NewTestWriter(t)})
			h := &Handler{
				logger: &logger,
				as:     mockStore,
			}

			ctx, cancel := tt.setupContext()
			defer cancel()

			resp, err := h.GetTags(ctx, &pb.Empty{})

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				if tt.name == "Large Number of Tags" {
					assert.Equal(t, 10000, len(resp.Tags))
				} else {
					assert.Equal(t, tt.expectedTags, resp.Tags)
				}
			}

			mockStore.AssertExpectations(t)
		})
	}
}
