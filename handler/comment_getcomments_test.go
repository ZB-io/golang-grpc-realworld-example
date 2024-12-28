package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock interfaces

type MockArticleService struct {
	mock.Mock
}

func (m *MockArticleService) GetByID(id uint) (*model.Article, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleService) GetComments(article *model.Article) ([]model.Comment, error) {
	args := m.Called(article)
	return args.Get(0).([]model.Comment), args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) IsFollowing(user *model.User, author *model.User) (bool, error) {
	args := m.Called(user, author)
	return args.Bool(0), args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info() *mock.Mock {
	m.Called()
	return &m.Mock
}

func (m *MockLogger) Error() *mock.Mock {
	m.Called()
	return &m.Mock
}

// Test Function
func TestGetComments(t *testing.T) {
	mockAS := new(MockArticleService)
	mockUS := new(MockUserService)
	mockLogger := new(MockLogger)

	handler := &Handler{
		as:     mockAS,
		us:     mockUS,
		logger: mockLogger,
	}

	// Table-driven tests
	tests := []struct {
		name          string
		setupMocks    func()
		slug          string
		expectedErr   error
		expectedCodes codes.Code
	}{
		{
			name: "Valid Article ID with Comments",
			setupMocks: func() {
				mockAS.On("GetByID", uint(1)).Return(&model.Article{}, nil)
				mockAS.On("GetComments", mock.Anything).Return([]model.Comment{
					{
						Author: model.User{},
					},
				}, nil)
				mockUS.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)
			},
			slug:          "1",
			expectedErr:   nil,
			expectedCodes: codes.OK,
		},
		{
			name:       "Invalid Article Slug Conversion",
			setupMocks: func() {},
			slug:       "invalid",
			expectedErr: status.Error(codes.InvalidArgument, "invalid article id"),
			expectedCodes: codes.InvalidArgument,
		},
		{
			name: "Article Not Found",
			setupMocks: func() {
				mockAS.On("GetByID", uint(2)).Return(nil, errors.New("article not found"))
			},
			slug:          "2",
			expectedErr:   status.Error(codes.InvalidArgument, "invalid article id"),
			expectedCodes: codes.InvalidArgument,
		},
		{
			name: "No Comments Available",
			setupMocks: func() {
				mockAS.On("GetByID", uint(3)).Return(&model.Article{}, nil)
				mockAS.On("GetComments", mock.Anything).Return([]model.Comment{}, nil)
			},
			slug:          "3",
			expectedErr:   nil,
			expectedCodes: codes.OK,
		},
		{
			name: "Unauthorized User Access",
			setupMocks: func() {
				mockAS.On("GetByID", uint(4)).Return(&model.Article{}, nil)
				mockAS.On("GetComments", mock.Anything).Return([]model.Comment{
					{
						Author: model.User{},
					},
				}, nil)
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, errors.New("not authenticated")
				}
			},
			slug:          "4",
			expectedErr:   status.Error(codes.Unauthenticated, "no user authentication"),
			expectedCodes: codes.Unauthenticated,
		},
		{
			name: "Error Retrieving Comments",
			setupMocks: func() {
				mockAS.On("GetByID", uint(5)).Return(&model.Article{}, nil)
				mockAS.On("GetComments", mock.Anything).Return(nil, errors.New("comments retrieval failed"))
			},
			slug:          "5",
			expectedErr:   status.Error(codes.Aborted, "failed to get comments"),
			expectedCodes: codes.Aborted,
		},
		{
			name: "Failing to Retrieve Following Status",
			setupMocks: func() {
				mockAS.On("GetByID", uint(6)).Return(&model.Article{}, nil)
				mockAS.On("GetComments", mock.Anything).Return([]model.Comment{
					{
						Author: model.User{},
					},
				}, nil)
				mockUS.On("IsFollowing", mock.Anything, mock.Anything).Return(false, errors.New("user following retrieval failed"))
			},
			slug:          "6",
			expectedErr:   status.Error(codes.NotFound, "internal server error"),
			expectedCodes: codes.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			req := &pb.GetCommentsRequest{Slug: tc.slug}
			resp, err := handler.GetComments(context.Background(), req)

			if tc.expectedErr != nil {
				assert.Nil(t, resp)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.expectedCodes, st.Code())
				assert.EqualError(t, err, fmt.Sprintf("%v", tc.expectedErr))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			t.Log("Passed test case:", tc.name)
		})
	}
}
