package handler

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlerUnfavoriteArticle(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func() context.Context
		setupMocks     func(sqlmock.Sqlmock)
		slug           string
		expectedResult *proto.ArticleResponse
		expectedError  error
	}{
		{
			name: "Successfully Unfavoriting an Article",
			setupContext: func() context.Context {
				return auth.WithToken(context.Background(), "Bearer validToken")
			},
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE id=\\?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("SELECT \\* FROM `articles` WHERE id=\\?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "description", "body", "favorites_count"}).
						AddRow(1, "1", "Test Title", "Test Description", "Test Body", 1))
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `favorites` WHERE `article_id` = \\? AND `user_id` = \\?").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE `articles` SET `favorites_count` = `favorites_count` - \\? WHERE `id` = \\?").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery("SELECT \\* FROM `follows` WHERE \\`from_user_id\\` = \\? AND \\`to_user_id\\` = \\?").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			slug: "1",
			expectedResult: &proto.ArticleResponse{
				Article: &proto.Article{
					Slug:           "1",
					Title:          "Test Title",
					Description:    "Test Description",
					Body:           "Test Body",
					FavoritesCount: 0,
					Favorited:      false,
					Author:         &proto.Profile{Username: "testuser", Following: false},
				},
			},
			expectedError: nil,
		},
		{
			name: "Handling an Unauthenticated User",
			setupContext: func() context.Context {
				return context.Background()
			},
			setupMocks:    func(mock sqlmock.Sqlmock) {},
			slug:          "1",
			expectedResult: nil,
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		// TODO: Add more test cases for other scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.setupMocks(mock)

			logger := zerolog.New(zerolog.ConsoleWriter{Out: zerolog.Nop()})
			us := &store.UserStore{DB: db}
			as := &store.ArticleStore{DB: db}
			h := &Handler{logger: &logger, us: us, as: as}

			ctx := tt.setupContext()

			req := &proto.UnfavoriteArticleRequest{Slug: tt.slug}
			resp, err := h.UnfavoriteArticle(ctx, req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, resp)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

