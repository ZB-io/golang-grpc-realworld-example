package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

func TestHandlerFavoriteArticle(t *testing.T) {
	// Initialize a zerolog logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: t})

	// Setting up mock user and article stores
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database connection: %s", err)
	}
	defer db.Close()

	userStore := &store.UserStore{DB: db}
	articleStore := &store.ArticleStore{DB: db}

	handler := &Handler{
		logger: &logger,
		us:     userStore,
		as:     articleStore,
	}

	// A structure for the test cases
	type testCase struct {
		desc        string
		mockSetups  func()
		req         *pb.FavoriteArticleRequest
		expectedErr error
		verify      func(result *pb.ArticleResponse, err error)
	}

	// Define test cases
	tests := []testCase{
		{
			desc: "Successful Article Favorite",
			mockSetups: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("SELECT (.+) FROM articles WHERE id = ?").
					WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(2, "test article"))
				mock.ExpectExec("INSERT INTO favorited_users").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			req: &pb.FavoriteArticleRequest{Slug: "2"},
			verify: func(result *pb.ArticleResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "test article", result.Article.Title)
			},
		},
		{
			desc: "Unauthenticated User",
			mockSetups: func() {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, errors.New("unauthenticated")
				}
			},
			req: &pb.FavoriteArticleRequest{Slug: "3"},
			verify: func(result *pb.ArticleResponse, err error) {
				assert.Nil(t, result)
				assert.EqualError(t, err, status.Error(codes.Unauthenticated, "unauthenticated").Error())
			},
		},
		{
			desc: "Non-existent User",
			mockSetups: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs(1).WillReturnError(fmt.Errorf("record not found"))
			},
			req: &pb.FavoriteArticleRequest{Slug: "4"},
			verify: func(result *pb.ArticleResponse, err error) {
				assert.Nil(t, result)
				assert.EqualError(t, err, status.Error(codes.NotFound, "not user found").Error())
			},
		},
		{
			desc: "Invalid Slug Format",
			mockSetups: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
			},
			req: &pb.FavoriteArticleRequest{Slug: "abc"},
			verify: func(result *pb.ArticleResponse, err error) {
				assert.Nil(t, result)
				assert.EqualError(t, err, status.Error(codes.InvalidArgument, "invalid article id").Error())
			},
		},
		{
			desc: "Article Not Found",
			mockSetups: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("SELECT (.+) FROM articles WHERE id = ?").
					WithArgs(5).WillReturnError(fmt.Errorf("record not found"))
			},
			req: &pb.FavoriteArticleRequest{Slug: "5"},
			verify: func(result *pb.ArticleResponse, err error) {
				assert.Nil(t, result)
				assert.EqualError(t, err, status.Error(codes.InvalidArgument, "invalid article id").Error())
			},
		},
		{
			desc: "Failure to Add Favorite",
			mockSetups: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("SELECT (.+) FROM articles WHERE id = ?").
					WithArgs(6).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(6, "test article"))
				mock.ExpectExec("INSERT INTO favorited_users").
					WillReturnError(fmt.Errorf("unable to add favorite"))
			},
			req: &pb.FavoriteArticleRequest{Slug: "6"},
			verify: func(result *pb.ArticleResponse, err error) {
				assert.Nil(t, result)
				assert.EqualError(t, err, status.Error(codes.InvalidArgument, "failed to add favorite").Error())
			},
		},
		{
			desc: "Failed to Retrieve Following Status",
			mockSetups: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "testuser"))
				mock.ExpectQuery("SELECT (.+) FROM articles WHERE id = ?").
					WithArgs(7).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(7, "test article"))
				mock.ExpectExec("INSERT INTO favorited_users").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery("SELECT (.+) FROM follows").
					WillReturnError(fmt.Errorf("internal error"))
			},
			req: &pb.FavoriteArticleRequest{Slug: "7"},
			verify: func(result *pb.ArticleResponse, err error) {
				assert.Nil(t, result)
				assert.EqualError(t, err, status.Error(codes.NotFound, "internal server error").Error())
			},
		},
	}

	// Run each test scenario
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockSetups()
			// Call the function
			res, err := handler.FavoriteArticle(context.Background(), tc.req)
			tc.verify(res, err)
			// Validate mock expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
