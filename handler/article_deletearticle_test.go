package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
)



type Logger struct{}


func (l *Logger) Err(err error) *Logger { return l }
func (l *Logger) Error() *Logger { return l }
func (l *Logger) Info() *Logger { return l }
func (l *Logger) Interface(key string, val interface{}) *Logger { return l }
func (l *Logger) Msg(msg string) {}
func TestHandlerDeleteArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := auth.NewMockAuthenticator(ctrl)
	mockUserStore := model.NewMockUserStore(ctrl)
	mockArticleStore := model.NewMockArticleStore(ctrl)
	h := &Handler{
		logger: &Logger{},
		us:     mockUserStore,
		as:     mockArticleStore,
	}

	type testCase struct {
		desc           string
		setupMocks     func()
		req            *pb.DeleteArticleRequest
		expectedCode   codes.Code
		expectedResult *pb.Empty
		expectedError  error
	}

	testCases := []testCase{
		{
			desc: "Valid Article Deletion",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserStore.EXPECT().GetByID(1).Return(&model.User{ID: 1}, nil)
				mockArticleStore.EXPECT().GetByID(uint(1)).Return(&model.Article{ID: 1, Author: &model.User{ID: 1}}, nil)
				mockArticleStore.EXPECT().Delete(gomock.Any()).Return(nil)
			},
			req:            &pb.DeleteArticleRequest{Slug: "1"},
			expectedCode:   codes.OK,
			expectedResult: &pb.Empty{},
		},
		{
			desc: "Unauthenticated User",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(0, errors.New("unauthenticated"))
			},
			req:           &pb.DeleteArticleRequest{Slug: "1"},
			expectedCode:  codes.Unauthenticated,
			expectedError: status.Errorf(codes.Unauthenticated, "unauthenticated"),
		},
		{
			desc: "User Not Found in Database",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserStore.EXPECT().GetByID(1).Return(nil, errors.New("user not found"))
			},
			req:           &pb.DeleteArticleRequest{Slug: "1"},
			expectedCode:  codes.NotFound,
			expectedError: status.Error(codes.NotFound, "not user found"),
		},
		{
			desc:         "Invalid Article Slug Format",
			setupMocks:   func() {},
			req:          &pb.DeleteArticleRequest{Slug: "invalid-slug"},
			expectedCode: codes.InvalidArgument,
			expectedError: status.Error(
				codes.InvalidArgument, "invalid article id",
			),
		},
		{
			desc: "Article Not Found",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserStore.EXPECT().GetByID(1).Return(&model.User{ID: 1}, nil)
				mockArticleStore.EXPECT().GetByID(uint(1)).Return(nil, errors.New("article not found"))
			},
			req:           &pb.DeleteArticleRequest{Slug: "1"},
			expectedCode:  codes.InvalidArgument,
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			desc: "Unauthorized Article Deletion Attempt",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserStore.EXPECT().GetByID(1).Return(&model.User{ID: 1}, nil)
				mockArticleStore.EXPECT().GetByID(uint(1)).Return(&model.Article{ID: 1, Author: &model.User{ID: 2}}, nil)
			},
			req:           &pb.DeleteArticleRequest{Slug: "1"},
			expectedCode:  codes.PermissionDenied,
			expectedError: status.Errorf(codes.PermissionDenied, "forbidden"),
		},
		{
			desc: "Article Deletion Fails",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1, nil)
				mockUserStore.EXPECT().GetByID(1).Return(&model.User{ID: 1}, nil)
				mockArticleStore.EXPECT().GetByID(uint(1)).Return(&model.Article{ID: 1, Author: &model.User{ID: 1}}, nil)
				mockArticleStore.EXPECT().Delete(gomock.Any()).Return(errors.New("deletion error"))
			},
			req:           &pb.DeleteArticleRequest{Slug: "1"},
			expectedCode:  codes.Internal,
			expectedError: status.Errorf(codes.Internal, "failed to delete article"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.setupMocks()

			res, err := h.DeleteArticle(context.Background(), tc.req)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedCode, status.Code(err), "Expected error code did not match")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, res)
			}
		})
	}
}


