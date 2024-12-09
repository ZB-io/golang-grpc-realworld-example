package handler

import (
	"context"
	"testing"
	"time"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockArticleStore is a mock implementation of an ArticleStore
type MockArticleStore struct {
	mock.Mock
}

func (m *MockArticleStore) GetTags() ([]Tag, error) {
	args := m.Called()
	return args.Get(0).([]Tag), args.Error(1)
}

// MockLogger is a placeholder to capture logs.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info() *MockLogger {
	m.Called()
	return m
}

func (m *MockLogger) Err(err error) *MockLogger {
	m.Called(err)
	return m
}

func (m *MockLogger) Interface(key string, value interface{}) *MockLogger {
	m.Called(key, value)
	return m
}

func (m *MockLogger) Msg(msg string) {
	m.Called(msg)
}

func TestGetTags(t *testing.T) {
	// Table-driven tests
	tests := []struct {
		name          string
		mockTags      []Tag
		mockError     error
		expectedTags  []string
		expectedError error
		cancelContext bool
	}{
		{
			name: "Successful Retrieval of Tags",
			mockTags: []Tag{
				{Name: "golang"},
				{Name: "grpc"},
			},
			mockError:     nil,
			expectedTags:  []string{"golang", "grpc"},
			expectedError: nil,
		},
		{
			name:          "Error Handling when GetTags Returns Error",
			mockTags:      nil,
			mockError:     sql.ErrConnDone, // Simulate a failure
			expectedTags:  nil,
			expectedError: status.Error(codes.Aborted, "internal server error"),
		},
		{
			name:          "Empty Tags List Returned",
			mockTags:      []Tag{},
			mockError:     nil,
			expectedTags:  []string{},
			expectedError: nil,
		},
		{
			name:          "Context Cancellation Handling",
			mockTags:      nil,
			mockError:     nil,
			expectedTags:  nil,
			expectedError: context.Canceled,
			cancelContext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockArticleStore := new(MockArticleStore)
			mockLogger := new(MockLogger)

			mockLogger.On("Info").Return(mockLogger)
			mockLogger.On("Interface", mock.Anything, mock.Anything).Return(mockLogger)
			mockLogger.On("Msg", mock.Anything)

			if tt.mockError != nil {
				mockLogger.On("Err", mock.Anything).Return(mockLogger)
			}

			mockArticleStore.On("GetTags").Return(tt.mockTags, tt.mockError)

			handler := &Handler{
				as:     mockArticleStore,
				logger: mockLogger,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			if tt.cancelContext {
				cancel()
			}

			resp, err := handler.GetTags(ctx, &pb.Empty{})

			if tt.expectedError != nil {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedTags, resp.Tags)
			}

			// Validate logger was called as expected
			mockLogger.AssertNumberOfCalls(t, "Info", 1)
		})
	}
}
