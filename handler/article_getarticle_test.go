package handler

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock the dependencies for testing
func setupMock(t *testing.T) (*Handler, sqlmock.Sqlmock, context.Context) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Could not initialize DB mock: %v", err)
	}

	logger := zerolog.New(nil) // Customize logger setup as needed
	us := &store.UserStore{db}
	as := &store.ArticleStore{db}
	handler := &Handler{logger: &logger, us: us, as: as}

	ctx := context.Background()
	return handler, mock, ctx
}

func TestHandlerGetArticle(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(mock sqlmock.Sqlmock, ctx context.Context) // For setting the desired mock state
		request       *pb.GetArticleRequest
		expectedError codes.Code
	}{
		{
			name: "Successfully Retrieve Article with Valid Slug and Authenticated User",
			setupMocks: func(mock sqlmock.Sqlmock, ctx context.Context) {
				// Mock auth.GetUserID to return a valid user ID
				userID := 1 // Change as needed for testing
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return uint(userID), nil
				}

				// Mock the ArticleStore GetByID, IsFavorited, and UserStore IsFollowing methods
				mock.ExpectQuery(`SELECT (.+) FROM "articles"`).
					WithArgs(articleID).
					WillReturnRows(sqlmock.NewRows([]string{"id", ...})).
					WithRow(1, "Title", ...) // TODO: Define expected columns and data

				mock.ExpectQuery(`SELECT COUNT\(.+\) FROM "favorite_articles"`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				mock.ExpectQuery(`SELECT COUNT\(.+\) FROM "follows"`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			request:       &pb.GetArticleRequest{Slug: "1"},
			expectedError: codes.OK,
		},
		{
			name: "Successfully Retrieve Article with Valid Slug and Anonymous User",
			setupMocks: func(mock sqlmock.Sqlmock, ctx context.Context) {
				// Mock auth.GetUserID to return an error indicating unauthenticated user
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 0, status.Error(codes.Unauthenticated, "unauthenticated")
				}

				mock.ExpectQuery(`SELECT (.+) FROM "articles"`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", ...}).
					WithRow(1, "Title", ...)) // TODO: Define expected columns and data
			},
			request:       &pb.GetArticleRequest{Slug: "1"},
			expectedError: codes.OK,
		},
		{
			name: "Handle Invalid Slug Format",
			setupMocks: func(mock sqlmock.Sqlmock, ctx context.Context) {
				// No DB mocking necessary
			},
			request:       &pb.GetArticleRequest{Slug: "invalid_slug"},
			expectedError: codes.InvalidArgument,
		},
		{
			name: "Handle Article Not Found",
			setupMocks: func(mock sqlmock.Sqlmock, ctx context.Context) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles"`).
					WithArgs(999).
					WillReturnError(fmt.Errorf("record not found")) // NO SQL row
			},
			request:       &pb.GetArticleRequest{Slug: "999"},
			expectedError: codes.InvalidArgument,
		},
		{
			name: "Handle Non-existent User ID from Token",
			setupMocks: func(mock sqlmock.Sqlmock, ctx context.Context) {
				userID := 999 // Dummy non-existent user
				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return uint(userID), nil
				}

				mock.ExpectQuery(`SELECT (.+) FROM "users"`).
					WithArgs(userID).
					WillReturnError(fmt.Errorf("user not found"))
			},
			request:       &pb.GetArticleRequest{Slug: "1"},
			expectedError: codes.NotFound,
		},
		{
			name: "Handle Internal Server Error on Favorited Check",
			setupMocks: func(mock sqlmock.Sqlmock, ctx context.Context) {
				// Mocking same as first test step due to get article success
				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", ...}).
					AddRow(1, "Sample Article", ...)) // Adjust columns/results as per your schema 

				mock.ExpectQuery(`SELECT count\(\*\) FROM "favorite_articles"`).
					WillReturnError(fmt.Errorf("internal server error"))

				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return uint(1), nil // setting valid user ID
				}
			},
			request:       &pb.GetArticleRequest{Slug: "1"},
			expectedError: codes.Aborted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mock, ctx := setupMock(t)
			if tt.setupMocks != nil {
				tt.setupMocks(mock, ctx)
			}

			resp, err := handler.GetArticle(ctx, tt.request)

			if tt.expectedError != codes.OK {
				if status.Code(err) != tt.expectedError {
					t.Errorf("Expected error code %v, got %v", tt.expectedError, status.Code(err))
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if resp == nil || resp.Article == nil {
					t.Error("Expected non-nil ArticleResponse, got nil")
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all SQL expectations met: %v", err)
			}
		})
	}
}

// Note: This test code makes several assumptions on the function behavior, mock setup, and test compatibility.
