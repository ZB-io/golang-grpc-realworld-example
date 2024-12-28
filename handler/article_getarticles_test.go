package handler

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Ensure that this struct matches whatever dependencies are needed in GetArticles
type Handler struct {
	as ArticleService
	us UserService
}

type ArticleService interface {
	GetArticles(tag, author string, favoritedBy *model.User, limit, offset int32) ([]model.Article, error)
	IsFavorited(article *model.Article, user *model.User) (bool, error)
}

type UserService interface {
	GetByUsername(username string) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	IsFollowing(follower, followee *model.User) (bool, error)
}

func TestHandlerGetArticles(t *testing.T) {
	type testCase struct {
		description      string
		setup            func(*sqlmock.Sqlmock)
		request          *pb.GetArticlesRequest
		expectedArticles int
		expectedError    error
	}

	tests := []testCase{
		{
			description: "Scenario 1: Retrieve Articles with Default Limit",
			setup: func(mock *sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2).AddRow(3)) // Add more rows to simulate > 20 articles
			},
			request:          &pb.GetArticlesRequest{Limit: 0},
			expectedArticles: 3, // Adjusted based on mock setup
			expectedError:    nil,
		},
		{
			description: "Scenario 2: Retrieve Articles Favorited by a Specific User",
			setup: func(mock *sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(4))
			},
			request:          &pb.GetArticlesRequest{Favorited: "testuser"},
			expectedArticles: 1,
			expectedError:    nil,
		},
		{
			description: "Scenario 3: Handle Non-Existent Favorite User Gracefully",
			setup: func(mock *sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WillReturnError(sqlmock.ErrNoRows)
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			request:          &pb.GetArticlesRequest{Favorited: "nonExistentUser"},
			expectedArticles: 0,
			expectedError:    nil,
		},
		{
			description: "Scenario 4: Retrieve Articles of a Specific Author",
			setup: func(mock *sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE author_id").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))
			},
			request:          &pb.GetArticlesRequest{Author: "specificAuthor"},
			expectedArticles: 1,
			expectedError:    nil,
		},
		{
			description: "Scenario 5: User Not Found in Context",
			setup: func(mock *sqlmock.Sqlmock) {
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, status.Error(codes.NotFound, "user not found")
				}
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			},
			request:          &pb.GetArticlesRequest{},
			expectedArticles: 0,
			expectedError:    status.Error(codes.NotFound, "user not found"),
		},
		{
			description: "Scenario 6: Internal Server Error on Article Retrieval",
			setup: func(mock *sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM articles WHERE").WillReturnError(status.Error(codes.Aborted, "internal server error"))
			},
			request:          &pb.GetArticlesRequest{},
			expectedArticles: 0,
			expectedError:    status.Error(codes.Aborted, "internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Setup Mock and Context
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			h := &Handler{
				as: /* Populate with mocked ArticleService */,
				us: /* Populate with mocked UserService */,
			}

			// Setup database expectations
			tt.setup(&mock)

			// Act: Call the GetArticles function
			resp, err := h.GetArticles(context.Background(), tt.request)

			// Assert expectations
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedArticles, len(resp.Articles))
			}
		})
	}
}

// Note: Mock implementations for as (ArticleService) and us (UserService) should be set up properly for execution.
