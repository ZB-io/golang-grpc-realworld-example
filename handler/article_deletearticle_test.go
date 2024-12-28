package handler

import (
	"context"
	"testing"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type ExpectedExec struct {
	queryBasedExpectation
	result driver.Result
	delay  time.Duration
}

type ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}

type ArticleStore struct {
	db *gorm.DB
}

type UserStore struct {
	db *gorm.DB
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}


type Logger struct {
	w       LevelWriter
	level   Level
	sampler Sampler
	context []byte
	hooks   []Hook
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestHandlerDeleteArticle(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zerolog.Nop()
	us := &store.UserStore{db}
	as := &store.ArticleStore{db}
	h := &Handler{logger: &logger, us: us, as: as}

	tests := []struct {
		name          string
		contextSetup  func() context.Context
		request       *proto.DeleteArticleRequest
		dbSetup       func()
		expectedError error
		verifyResult  func(t *testing.T, result *proto.Empty, err error)
	}{
		{
			name: "Valid Article Deletion",
			contextSetup: func() context.Context {
				ctx := context.TODO()
				authCtx := context.WithValue(ctx, "UserID", uint(1))
				return authCtx
			},
			request: &proto.DeleteArticleRequest{Slug: "1"},
			dbSetup: func() {
				mock.ExpectQuery("SELECT *").WithArgs(1).WillReturnRows(
					sqlmock.NewRows([]string{"ID", "Name"}).AddRow(1, "John Doe"))

				mock.ExpectQuery("SELECT *").WithArgs(1).WillReturnRows(
					sqlmock.NewRows([]string{"ID", "Slug", "AuthorID"}).AddRow(1, "1", 1))

				mock.ExpectExec("DELETE FROM").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
			verifyResult: func(t *testing.T, result *proto.Empty, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			},
		},
		{
			name: "Unauthenticated User",
			contextSetup: func() context.Context {
				return context.TODO()
			},
			request:       &proto.DeleteArticleRequest{Slug: "1"},
			dbSetup:       func() {},
			expectedError: status.Error(codes.Unauthenticated, "unauthenticated"),
			verifyResult: func(t *testing.T, result *proto.Empty, err error) {
				assert.Nil(t, result)
				assert.ErrorIs(t, err, status.Error(codes.Unauthenticated, "unauthenticated"))
			},
		},
		{
			name: "Invalid Article Slug Conversion",
			contextSetup: func() context.Context {
				ctx := context.TODO()
				authCtx := context.WithValue(ctx, "UserID", uint(1))
				return authCtx
			},
			request:       &proto.DeleteArticleRequest{Slug: "invalid-slug"},
			dbSetup:       func() {},
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
			verifyResult: func(t *testing.T, result *proto.Empty, err error) {
				assert.Nil(t, result)
				assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid article id"))
			},
		},
		{
			name: "Article Not Found",
			contextSetup: func() context.Context {
				ctx := context.TODO()
				authCtx := context.WithValue(ctx, "UserID", uint(1))
				return authCtx
			},
			request: &proto.DeleteArticleRequest{Slug: "999"},
			dbSetup: func() {
				mock.ExpectQuery("SELECT *").WithArgs(999).WillReturnError(
					errors.New("article not found"))
			},
			expectedError: status.Error(codes.InvalidArgument, "invalid article id"),
			verifyResult: func(t *testing.T, result *proto.Empty, err error) {
				assert.Nil(t, result)
				assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid article id"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dbSetup()

			ctx := tt.contextSetup()
			result, err := h.DeleteArticle(ctx, tt.request)

			tt.verifyResult(t, result, err)
			assert.Equal(t, tt.expectedError, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
