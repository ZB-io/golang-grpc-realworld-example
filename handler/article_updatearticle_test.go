package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestHandlerUpdateArticle tests the UpdateArticle function from the handler package.
func TestUpdateArticle(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateArticleRequest
	}

	type test struct {
		name           string
		args           args
		mockSetup      func(sqlmock.Sqlmock, *store.UserStore, *store.ArticleStore)
		expectedError  error
		expectedResult *pb.ArticleResponse
	}

	tests := []test{
		{
			name: "Successful Article Update by Author",
			args: args{
				ctx: context.WithValue(context.Background(), "user_id", 1), // mock the context with user ID 1
				req: &pb.UpdateArticleRequest{
					Article: &pb.UpdateArticleRequest_Article{
						Slug:        "1",
						Title:       "New Title",
						Description: "New Description",
						Body:        "New Body",
					},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock, us *store.UserStore, as *store.ArticleStore) {
				// Mock user retrieval
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
						AddRow(1, "author"))

				// Mock article retrieval
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id"}).
						AddRow(1, "Old Title", 1))

				// Mock article update
				mock.ExpectExec(`UPDATE "articles" SET "title"=\$1, "description"=\$2, "body"=\$3 WHERE id = \$4`).
					WithArgs("New Title", "New Description", "New Body", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock following check
				mock.ExpectQuery(`SELECT count\(\*\) FROM "follows" WHERE from_user_id = \$1 AND to_user_id = \$2`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedError: nil,
			expectedResult: &pb.ArticleResponse{
				Article: &pb.Article{
					Title: "New Title",
				},
			},
		},
		{
			name: "Unauthenticated User Attempt",
			args: args{
				ctx: context.Background(),
				req: &pb.UpdateArticleRequest{
					Article: &pb.UpdateArticleRequest_Article{
						Slug: "1",
					},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock, us *store.UserStore, as *store.ArticleStore) {
				// No specific database operations expected as the authentication should fail first
			},
			expectedError: status.Errorf(codes.Unauthenticated, "unauthenticated"),
			expectedResult: nil,
		},
		{
			name: "Attempt to Update Another User's Article",
			args: args{
				ctx: context.WithValue(context.Background(), "user_id", 2), // mock the context with user ID 2
				req: &pb.UpdateArticleRequest{
					Article: &pb.UpdateArticleRequest_Article{
						Slug: "1",
					},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock, us *store.UserStore, as *store.ArticleStore) {
				// Mock user retrieval
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1`).
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
						AddRow(2, "different_author"))

				// Mock article retrieval
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id"}).
						AddRow(1, "Old Title", 1))
			},
			expectedError: status.Errorf(codes.PermissionDenied, "forbidden"),
			expectedResult: nil,
		},
		{
			name: "Invalid Slug Conversion",
			args: args{
				ctx: context.WithValue(context.Background(), "user_id", 1),
				req: &pb.UpdateArticleRequest{
					Article: &pb.UpdateArticleRequest_Article{
						Slug: "abc",
					},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock, us *store.UserStore, as *store.ArticleStore) {
				// No database calls expected because slug conversion fails
			},
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
			expectedResult: nil,
		},
		{
			name: "Validation Failure on Updated Article Fields",
			args: args{
				ctx: context.WithValue(context.Background(), "user_id", 1),
				req: &pb.UpdateArticleRequest{
					Article: &pb.UpdateArticleRequest_Article{
						Slug:  "1",
						Title: "", // Invalid: Title is required
					},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock, us *store.UserStore, as *store.ArticleStore) {
				// Mock user retrieval
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
						AddRow(1, "author"))

				// Mock article retrieval
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id"}).
						AddRow(1, "Old Title", 1))
			},
			expectedError: status.Error(codes.InvalidArgument, "validation error: Field validation failed for 'Title', Field validation failed for 'Body', Field validation failed for 'Tags'"),
			expectedResult: nil,
		},
		{
			name: "Database Failure during Article Update",
			args: args{
				ctx: context.WithValue(context.Background(), "user_id", 1),
				req: &pb.UpdateArticleRequest{
					Article: &pb.UpdateArticleRequest_Article{
						Slug:  "1",
						Title: "New Title",
					},
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock, us *store.UserStore, as *store.ArticleStore) {
				// Mock user retrieval
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
						AddRow(1, "author"))

				// Mock article retrieval
				mock.ExpectQuery(`SELECT \* FROM "articles" WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author_id"}).
						AddRow(1, "Old Title", 1))

				// Mock article update failure
				mock.ExpectExec(`UPDATE "articles" SET "title"=\$1, "description"=\$2, "body"=\$3 WHERE id = \$4`).
					WithArgs("New Title", "New Description", "New Body", 1).
					WillReturnError(errors.New("database failure"))
			},
			expectedError: status.Error(codes.InvalidArgument, "internal server error"),
			expectedResult: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			sqlDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer sqlDB.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{})
			assert.NoError(t, err)

			userStore := store.NewUserStore(gormDB)
			articleStore := store.NewArticleStore(gormDB)

			logger := zerolog.New(nil)
			h := handler.Handler{
				Logger: &logger,
				US:     userStore,
				AS:     articleStore,
			}

			// Setup the mock expectations
			tc.mockSetup(mock, userStore, articleStore)

			// Act
			result, err := h.UpdateArticle(tc.args.ctx, tc.args.req)

			// Assert
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult.Article.Title, result.Article.Title)
			}
			t.Logf("Result: %+v, Error: %v", result, err)
		})
	}
}
