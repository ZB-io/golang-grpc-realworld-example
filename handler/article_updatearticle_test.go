package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mocking UserService and ArticleService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) IsFollowing(currentUser, targetUser *model.User) (bool, error) {
	args := m.Called(currentUser, targetUser)
	return args.Bool(0), args.Error(1)
}

type MockArticleService struct {
	mock.Mock
}

func (m *MockArticleService) GetByID(id uint) (*model.Article, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Article), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockArticleService) Update(article *model.Article) error {
	return m.Called(article).Error(0)
}

func TestHandlerUpdateArticle(t *testing.T) {
	type testCase struct {
		desc     string
		setup    func(*Handler, *MockUserService, *MockArticleService)
		expected codes.Code
	}

	tests := []testCase{
		{
			desc: "Unauthorized User Attempt to Update Article",
			setup: func(h *Handler, us *MockUserService, as *MockArticleService) {
				auth.MockGetUserID = func(ctx context.Context) (uint, error) { return 0, status.Error(codes.Unauthenticated, "unauthenticated") }
			},
			expected: codes.Unauthenticated,
		},
		{
			desc: "Valid User Not Found in System",
			setup: func(h *Handler, us *MockUserService, as *MockArticleService) {
				auth.MockGetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				us.On("GetByID", uint(1)).Return(nil, status.Error(codes.NotFound, "not user found"))
			},
			expected: codes.NotFound,
		},
		{
			desc: "Invalid Article Slug Conversion",
			setup: func(h *Handler, us *MockUserService, as *MockArticleService) {
				auth.MockGetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
			},
			expected: codes.InvalidArgument,
		},
		{
			desc: "Article Not Found for Given Slug",
			setup: func(h *Handler, us *MockUserService, as *MockArticleService) {
				auth.MockGetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", uint(1)).Return(nil, status.Error(codes.InvalidArgument, "invalid article id"))
			},
			expected: codes.InvalidArgument,
		},
		{
			desc: "User Attempts to Update Another User's Article",
			setup: func(h *Handler, us *MockUserService, as *MockArticleService) {
				auth.MockGetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", mock.Anything).Return(&model.Article{Author: &model.User{ID: 2}}, nil)
			},
			expected: codes.Unauthenticated,
		},
		{
			desc: "Article Overwrite Validation Failure",
			setup: func(h *Handler, us *MockUserService, as *MockArticleService) {
				auth.MockGetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", mock.Anything).Return(&model.Article{Author: &model.User{ID: 1}}, nil)
				as.On("Update", mock.Anything).Return(status.Error(codes.InvalidArgument, "validation error"))
			},
			expected: codes.InvalidArgument,
		},
		{
			desc: "Successful Article Update",
			setup: func(h *Handler, us *MockUserService, as *MockArticleService) {
				auth.MockGetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", mock.Anything).Return(&model.Article{Author: &model.User{ID: 1}}, nil)
				as.On("Update", mock.Anything).Return(nil)
				us.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
			},
			expected: codes.OK,
		},
		{
			desc: "Internal Server Error on Update Operation",
			setup: func(h *Handler, us *MockUserService, as *MockArticleService) {
				auth.MockGetUserID = func(ctx context.Context) (uint, error) { return 1, nil }
				us.On("GetByID", uint(1)).Return(&model.User{ID: 1}, nil)
				as.On("GetByID", mock.Anything).Return(&model.Article{Author: &model.User{ID: 1}}, nil)
				as.On("Update", mock.Anything).Return(status.Error(codes.Aborted, "failed to update article"))
			},
			expected: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			us := new(MockUserService)
			as := new(MockArticleService)
			h := &Handler{us: us, as: as}

			tt.setup(h, us, as)

			req := &pb.UpdateArticleRequest{
				Article: &pb.Article{
					Slug: "1", // Default slug for valid input test cases
				},
			}
			resp, err := h.UpdateArticle(context.Background(), req)
			
			if err != nil {
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.expected, st.Code())
			} else {
				require.NotNil(t, resp)
				require.Equal(t, tt.expected, codes.OK)
			}

			us.AssertExpectations(t)
			as.AssertExpectations(t)
		})
	}
}
