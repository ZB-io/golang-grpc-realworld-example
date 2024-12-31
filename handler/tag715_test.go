package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		expectedTags   []string
		expectedError  error
		contextTimeout time.Duration
	}{
		{
			name: "Successful Retrieval of Tags",
			setupMock: func(m *MockArticleStore) {
				m.On("GetTags").Return([]model.Tag{{Name: "tag1"}, {Name: "tag2"}}, nil)
			},
			expectedTags:  []string{"tag1", "tag2"},
			expectedError: nil,
		},
		{
			name: "Empty Tag List",
			setupMock: func(m *MockArticleStore) {
				m.On("GetTags").Return([]model.Tag{}, nil)
			},
			expectedTags:  []string{},
			expectedError: nil,
		},
		{
			name: "Error from ArticleStore",
			setupMock: func(m *MockArticleStore) {
				m.On("GetTags").Return([]model.Tag{}, errors.New("store error"))
			},
			expectedTags:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
		},
		{
			name: "Context Cancellation",
			setupMock: func(m *MockArticleStore) {
				m.On("GetTags").After(100*time.Millisecond).Return([]model.Tag{}, nil)
			},
			expectedTags:   nil,
			expectedError:  context.Canceled,
			contextTimeout: 50 * time.Millisecond,
		},
		{
			name: "Large Number of Tags",
			setupMock: func(m *MockArticleStore) {
				tags := make([]model.Tag, 10000)
				for i := range tags {
					tags[i] = model.Tag{Name: fmt.Sprintf("tag%d", i)}
				}
				m.On("GetTags").Return(tags, nil)
			},
			expectedTags:  nil,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockArticleStore)
			tt.setupMock(mockStore)

			logger := zerolog.New(zerolog.NewConsoleWriter())
			h := &Handler{
				logger: &logger,
				as:     mockStore,
			}

			ctx := context.Background()
			if tt.contextTimeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.contextTimeout)
				defer cancel()
			}

			resp, err := h.GetTags(ctx, &pb.Empty{})

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
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
