package handler

import (
	"context"
	"errors"
	"testing"
	"time"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
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

FUNCTION_DEF=func (h *Handler) GetTags(ctx context.Context, req *pb.Empty) (*pb.TagsResponse, error) 

 */
func (m *MockArticleStore) GetTags() ([]model.Tag, error) {
	args := m.Called()
	return args.Get(0).([]model.Tag), args.Error(1)
}

func TestHandlerGetTags(t *testing.T) {

	logger := zerolog.New(zerolog.NewTestWriter(t))

	tests := []struct {
		name          string
		setupMock     func(*MockArticleStore)
		ctx           context.Context
		expectedTags  []string
		expectedError error
		validateError func(*testing.T, error)
		setupContext  func() (context.Context, context.CancelFunc)
	}{
		{
			name: "Successful Retrieval of Tags",
			setupMock: func(mas *MockArticleStore) {
				mas.On("GetTags").Return([]model.Tag{
					{Name: "golang"},
					{Name: "testing"},
				}, nil)
			},
			ctx:           context.Background(),
			expectedTags:  []string{"golang", "testing"},
			expectedError: nil,
			validateError: nil,
		},
		{
			name: "Empty Tags List",
			setupMock: func(mas *MockArticleStore) {
				mas.On("GetTags").Return([]model.Tag{}, nil)
			},
			ctx:           context.Background(),
			expectedTags:  []string{},
			expectedError: nil,
			validateError: nil,
		},
		{
			name: "Database Error Handling",
			setupMock: func(mas *MockArticleStore) {
				mas.On("GetTags").Return([]model.Tag{}, errors.New("database error"))
			},
			ctx:           context.Background(),
			expectedTags:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
			validateError: func(t *testing.T, err error) {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.Aborted, st.Code())
			},
		},
		{
			name: "Context Cancellation",
			setupMock: func(mas *MockArticleStore) {
				mas.On("GetTags").Return([]model.Tag{}, context.Canceled)
			},
			setupContext: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
				return ctx, cancel
			},
			expectedTags:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
			validateError: func(t *testing.T, err error) {
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.Aborted, st.Code())
			},
		},
		{
			name: "Duplicate Tag Handling",
			setupMock: func(mas *MockArticleStore) {
				mas.On("GetTags").Return([]model.Tag{
					{Name: "golang"},
					{Name: "golang"},
					{Name: "testing"},
				}, nil)
			},
			ctx:           context.Background(),
			expectedTags:  []string{"golang", "golang", "testing"},
			expectedError: nil,
			validateError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockStore := new(MockArticleStore)
			tt.setupMock(mockStore)

			h := &Handler{
				logger: &logger,
				as:     mockStore,
			}

			ctx := tt.ctx
			if tt.setupContext != nil {
				var cancel context.CancelFunc
				ctx, cancel = tt.setupContext()
				defer cancel()
			}

			response, err := h.GetTags(ctx, &pb.Empty{})

			if tt.validateError != nil {
				tt.validateError(t, err)
			} else if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.expectedTags, response.Tags)
			}

			mockStore.AssertExpectations(t)

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}

