package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteArticle(t *testing.T) {
	type test struct {
		name       string
		ctx        func() context.Context
		req        *pb.DeleteArticleRequest
		setupMocks func(userID uint, articleStore *store.ArticleStore, userStore *store.UserStore, sqlMock sqlmock.Sqlmock)
		wantErr    codes.Code
	}

	tests := []test{
		{
			name: "Valid Article Deletion",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), "user_id", uint(1))
			},
			req: &pb.DeleteArticleRequest{Slug: "42"},
			setupMocks: func(userID uint, articleStore *store.ArticleStore, userStore *store.UserStore, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery("^SELECT (.+) FROM users WHERE (.+)$").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				sqlMock.ExpectQuery("^SELECT (.+) FROM articles WHERE (.+)$").WillReturnRows(sqlmock.NewRows([]string{"id", "author_id"}).AddRow(42, 1))
				sqlMock.ExpectExec("^DELETE FROM articles WHERE (.+)$").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: codes.OK,
		},
		{
			name: "Unauthenticated User",
			ctx: func() context.Context {
				return context.Background()
			},
			req: &pb.DeleteArticleRequest{Slug: "42"},
			setupMocks: func(userID uint, articleStore *store.ArticleStore, userStore *store.UserStore, sqlMock sqlmock.Sqlmock) {
				// Override the default GetUserID function to simulate unauthentication
				auth.GetUserID = func(_ context.Context) (uint, error) {
					return 0, errors.New("unauthenticated")
				}
			},
			wantErr: codes.Unauthenticated,
		},
		{
			name: "Invalid Article Slug Conversion",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), "user_id", uint(1))
			},
			req: &pb.DeleteArticleRequest{Slug: "invalidSlug"},
			setupMocks: func(userID uint, articleStore *store.ArticleStore, userStore *store.UserStore, sqlMock sqlmock.Sqlmock) {
				// No specific mocks needed
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "Article Not Found",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), "user_id", uint(1))
			},
			req: &pb.DeleteArticleRequest{Slug: "999"},
			setupMocks: func(userID uint, articleStore *store.ArticleStore, userStore *store.UserStore, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery("^SELECT (.+) FROM articles WHERE (.+)$").WillReturnError(errors.New("record not found"))
			},
			wantErr: codes.NotFound,
		},
		{
			name: "Unauthorized Access to Another User's Article",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), "user_id", uint(1))
			},
			req: &pb.DeleteArticleRequest{Slug: "43"},
			setupMocks: func(userID uint, articleStore *store.ArticleStore, userStore *store.UserStore, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery("^SELECT (.+) FROM articles WHERE (.+)$").WillReturnRows(sqlmock.NewRows([]string{"id", "author_id"}).AddRow(43, 2))
			},
			wantErr: codes.PermissionDenied,
		},
		{
			name: "Article Deletion Failure in Store",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), "user_id", uint(1))
			},
			req: &pb.DeleteArticleRequest{Slug: "42"},
			setupMocks: func(userID uint, articleStore *store.ArticleStore, userStore *store.UserStore, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery("^SELECT (.+) FROM articles WHERE (.+)$").WillReturnRows(sqlmock.NewRows([]string{"id", "author_id"}).AddRow(42, 1))
				sqlMock.ExpectExec("^DELETE FROM articles WHERE (.+)$").WillReturnError(errors.New("delete failure"))
			},
			wantErr: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.New(nil)
			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			articleStore := &store.ArticleStore{DB: db}
			userStore := &store.UserStore{DB: db}

			tt.setupMocks(1, articleStore, userStore, sqlMock)

			h := &Handler{
				logger: &logger,
				us:     userStore,
				as:     articleStore,
			}

			_, err = h.DeleteArticle(tt.ctx(), tt.req)
			if status.Code(err) != tt.wantErr {
				t.Errorf("expected error code %v, got %v", tt.wantErr, status.Code(err))
			}
		})
	}
}
