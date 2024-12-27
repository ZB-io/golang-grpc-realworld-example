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
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/DATA-DOG/go-sqlmock"
)

/*
ROOST_METHOD_HASH=GetComments_265127fb6a
ROOST_METHOD_SIG_HASH=GetComments_20efd5abae
*/
func TestGetComments(t *testing.T) {
	mockAS := new(MockArticleService)
	mockUS := new(MockUserService)
	mockLogger := new(MockLogger)

	handler := &Handler{
		as:     mockAS,
		us:     mockUS,
		logger: mockLogger,
	}

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
			expectedErr:   status.Error(codes.Unauthenticated, "unauthenticated"),
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

/*
ROOST_METHOD_HASH=DeleteComment_452af2f984
ROOST_METHOD_SIG_HASH=DeleteComment_27615e7d69
*/
func TestDeleteComment(t *testing.T) {
	type testCase struct {
		name     string
		ctx      context.Context
		req      *pb.DeleteCommentRequest
		setup    func(mock sqlmock.Sqlmock, h *Handler)
		expected codes.Code
	}

	tests := []testCase{
		{
			name: "Successful Comment Deletion",
			ctx:  auth.SetUserID(context.Background(), 1),
			req: &pb.DeleteCommentRequest{
				Id:   "123",
				Slug: "123",
			},
			setup: func(mock sqlmock.Sqlmock, h *Handler) {
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE "id"=?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "email@example.com"))

				mock.ExpectQuery(`SELECT .* FROM "comments" WHERE "id"=?`).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "article_id"}).AddRow(123, 1, 123))

				mock.ExpectExec(`DELETE FROM "comments" WHERE "id"=?`).
					WithArgs(123).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expected: codes.OK,
		},
		{
			name: "Unauthenticated User",
			ctx:  context.Background(),
			req: &pb.DeleteCommentRequest{
				Id:   "123",
				Slug: "123",
			},
			setup: func(mock sqlmock.Sqlmock, h *Handler) {},
			expected: codes.Unauthenticated,
		},
		{
			name: "Invalid Comment ID Format",
			ctx:  auth.SetUserID(context.Background(), 1),
			req: &pb.DeleteCommentRequest{
				Id:   "abc",
				Slug: "123",
			},
			setup: func(mock sqlmock.Sqlmock, h *Handler) {
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE "id"=?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "email@example.com"))
			},
			expected: codes.InvalidArgument,
		},
		{
			name: "Comment Not Found",
			ctx:  auth.SetUserID(context.Background(), 1),
			req: &pb.DeleteCommentRequest{
				Id:   "123",
				Slug: "123",
			},
			setup: func(mock sqlmock.Sqlmock, h *Handler) {
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE "id"=?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "email@example.com"))

				mock.ExpectQuery(`SELECT .* FROM "comments" WHERE "id"=?`).
					WithArgs(123).
					WillReturnError(sqlmock.ErrNotFound)
			},
			expected: codes.InvalidArgument,
		},
		{
			name: "Comment Belongs to a Different User",
			ctx:  auth.SetUserID(context.Background(), 1),
			req: &pb.DeleteCommentRequest{
				Id:   "123",
				Slug: "123",
			},
			setup: func(mock sqlmock.Sqlmock, h *Handler) {
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE "id"=?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "email@example.com"))

				mock.ExpectQuery(`SELECT .* FROM "comments" WHERE "id"=?`).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "article_id"}).AddRow(123, 2, 123))
			},
			expected: codes.InvalidArgument,
		},
		{
			name: "Comment ID and Article Slug Mismatch",
			ctx:  auth.SetUserID(context.Background(), 1),
			req: &pb.DeleteCommentRequest{
				Id:   "123",
				Slug: "999",
			},
			setup: func(mock sqlmock.Sqlmock, h *Handler) {
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE "id"=?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "email@example.com"))

				mock.ExpectQuery(`SELECT .* FROM "comments" WHERE "id"=?`).
					WithArgs(123).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "article_id"}).AddRow(123, 1, 123))
			},
			expected: codes.InvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unexpected error opening a stub database connection: %s", err)
			}
			defer db.Close()

			handler := &Handler{}
			tc.setup(mock, handler)

			resp, err := handler.DeleteComment(tc.ctx, tc.req)
			if status.Code(err) != tc.expected {
				t.Errorf("expected error code %v, got %v", tc.expected, status.Code(err))
			}
			t.Logf("Test scenario '%s' resulted in response: %+v with error: %v", tc.name, resp, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=CreateComment_c4ccd62dc5
ROOST_METHOD_SIG_HASH=CreateComment_19a3ee5a3b
*/
func TestHandlerCreateComment(t *testing.T) {
	mockUserService := &MockUserService{}
	mockArticleService := &MockArticleService{}
	mockLogger := &MockLogger{}

	handler := &Handler{
		us:     mockUserService,
		as:     mockArticleService,
		logger: mockLogger,
	}

	tests := []struct {
		name   string
		setup  func()
		input  *pb.CreateCommentRequest
		verify func(t *testing.T, resp *pb.CommentResponse, err error)
	}{
		{
			name: "Scenario 1: Unauthorized User Cannot Create Comment",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, status.Error(codes.Unauthenticated, "unauthenticated")
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "This is a comment",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.Unauthenticated, status.Code(err))
				t.Log("Unauthorized user was correctly denied comment creation")
			},
		},
		{
			name: "Scenario 2: Non-existent User Cannot Create Comment",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return nil, status.Error(codes.NotFound, "user not found")
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "This is a comment",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.NotFound, status.Code(err))
				t.Log("Correctly identified non-existent user")
			},
		},
		{
			name: "Scenario 3: Invalid Article Slug Cannot Lead to Comment Creation",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "NaN",
				Comment: &pb.NewComment{
					Body: "This is a comment",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				t.Log("Invalid article slug was correctly rejected")
			},
		},
		{
			name: "Scenario 4: Non-existent Article Cannot Have a Comment",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
				mockArticleService.GetByIDMock = func(id uint) (*model.Article, error) {
					return nil, status.Error(codes.InvalidArgument, "invalid article id")
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "99",
				Comment: &pb.NewComment{
					Body: "Trying to comment on a non-existent article",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				t.Log("Non-existent article correctly identified")
			},
		},
		{
			name: "Scenario 5: Valid Comment on Existing Article",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
				mockArticleService.GetByIDMock = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1}, nil
				}
				mockArticleService.CreateCommentMock = func(comment *model.Comment) error {
					return nil
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "This is a valid comment",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.NoError(t, err)
				assert.NotNil(t, resp)
				t.Log("Comment successfully created")
			},
		},
		{
			name: "Scenario 6: Comment Fails Validation",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
				mockArticleService.GetByIDMock = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1}, nil
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				t.Log("Correctly failed on comment validation")
			},
		},
		{
			name: "Scenario 7: Database Error on Comment Creation",
			setup: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 42, nil
				}
				mockUserService.GetByIDMock = func(id uint) (*model.User, error) {
					return &model.User{ID: 42}, nil
				}
				mockArticleService.GetByIDMock = func(id uint) (*model.Article, error) {
					return &model.Article{ID: 1}, nil
				}
				mockArticleService.CreateCommentMock = func(comment *model.Comment) error {
					return status.Error(codes.Aborted, "database error")
				}
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.NewComment{
					Body: "This comment should test db error",
				},
			},
			verify: func(t *testing.T, resp *pb.CommentResponse, err error) {
				require.Error(t, err)
				assert.Equal(t, codes.Aborted, status.Code(err))
				t.Log("Database error correctly handled")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resp, err := handler.CreateComment(context.Background(), tt.input)
			tt.verify(t, resp, err)
		})
	}
}
