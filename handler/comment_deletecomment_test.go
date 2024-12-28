package handler

import (
	"context"
	"testing"
	"fmt"
	"strconv"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler represents the handler for gRPC requests.
type Handler struct {
	logger Logger
	us     UserService
	as     ArticleService
}

// Logger is an interface for structured logging.
type Logger interface {
	Info() Logger
	Error() Logger
	Msg(string)
	Msgf(string, ...interface{})
	Err(error) Logger
}

// UserService is an interface for handling user operations.
type UserService interface {
	GetByID(userID int) (*model.User, error)
}

// ArticleService is an interface for article operations.
type ArticleService interface {
	GetCommentByID(id uint) (*model.Comment, error)
	DeleteComment(comment *model.Comment) error
}

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
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "article_id"}).AddRow(123, 1, 123)) // Correct article
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

			handler := &Handler{
				// Mock necessary services here
				// logger, us, as needs setup
			}

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
