package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlerCreateComment(t *testing.T) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database connection: %s", err)
	}
	defer db.Close()

	userStore := &store.UserStore{DB: db}
	articleStore := &store.ArticleStore{DB: db}
	handler := &Handler{logger: &logger, us: userStore, as: articleStore}
	ctx := context.Background()

	tests := []struct {
		name          string
		setupMocks    func()
		input         *pb.CreateCommentRequest
		expectedResp  *pb.CommentResponse
		expectedError error
	}{
		{
			name: "Successful Comment Creation",
			setupMocks: func() {
				ctx = context.WithValue(ctx, auth.UserIDKey, uint(1))
				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"ID", "Username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("SELECT \\* FROM \"articles\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"ID", "Title"}).AddRow(1, "Test Article"))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"comments\" \\(body, user_id, article_id\\) VALUES \\(\\$1, \\$2, \\$3\\)").
					WithArgs("Great article!", 1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "Great article!",
				},
			},
			expectedResp: &pb.CommentResponse{Comment: &pb.Comment{
				Id:   "1",
				Body: "Great article!",
				Author: &pb.Profile{
					Username: "testuser",
				},
			}},
			expectedError: nil,
		},
		{
			name: "Unauthenticated User",
			setupMocks: func() {
				// No setup needed for unauthenticated context
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "Great article!",
				},
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "User Not Found",
			setupMocks: func() {
				ctx = context.WithValue(ctx, auth.UserIDKey, uint(1))
				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnError(fmt.Errorf("user not found"))
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "Nice article!",
				},
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.NotFound, "user not found"),
		},
		{
			name: "Invalid Article ID in Slug",
			setupMocks: func() {
				ctx = context.WithValue(ctx, auth.UserIDKey, uint(1))
				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"ID", "Username"}).AddRow(1, "testuser"))
			},
			input: &pb.CreateCommentRequest{
				Slug: "invalid_slug",
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "Nice article!",
				},
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Article Not Found",
			setupMocks: func() {
				ctx = context.WithValue(ctx, auth.UserIDKey, uint(1))
				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"ID", "Username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("SELECT \\* FROM \"articles\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnError(fmt.Errorf("article not found"))
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "Nice article!",
				},
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
		},
		{
			name: "Comment Validation Fails",
			setupMocks: func() {
				ctx = context.WithValue(ctx, auth.UserIDKey, uint(1))
				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"ID", "Username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("SELECT \\* FROM \"articles\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"ID", "Title"}).AddRow(1, "Test Article"))
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "",
				},
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.InvalidArgument, "validation error: cannot be blank"),
		},
		{
			name: "Comment Creation Aborted",
			setupMocks: func() {
				ctx = context.WithValue(ctx, auth.UserIDKey, uint(1))
				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"ID", "Username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("SELECT \\* FROM \"articles\" WHERE \\(id = \\$1\\)").WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"ID", "Title"}).AddRow(1, "Test Article"))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"comments\" \\(body, user_id, article_id\\) VALUES \\(\\$1, \\$2, \\$3\\)").
					WithArgs("Nice article!", 1, 1).
					WillReturnError(fmt.Errorf("comment creation failed"))
				mock.ExpectRollback()
			},
			input: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "Nice article!",
				},
			},
			expectedResp:  nil,
			expectedError: status.Error(codes.Aborted, "failed to create comment."),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			resp, err := handler.CreateComment(ctx, tc.input)
			assert.Equal(t, tc.expectedResp, resp)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test case %q executed successfully", tc.name)
		})
	}
}
