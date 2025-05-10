package handler

import (
	"context"
	"testing"
	"time"
	"os"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	// "github.com/raahii/golang-grpc-realworld-example/model" // Imported and not used
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlerGetFeedArticles(t *testing.T) {
	type testCase struct {
		name  string
		ctx   context.Context
		req   *pb.GetFeedArticlesRequest
		setup func(*store.UserStore, *store.ArticleStore)
		check func(*pb.ArticlesResponse, error)
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}) // os package used

	// Mocking database and necessary components
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storeUser := &store.UserStore{DB: db}     // Ensure these structs are properly defined elsewhere
	storeArticle := &store.ArticleStore{DB: db} // Ensure these structs are properly defined elsewhere

	handler := &Handler{
		logger: &logger, // Handler struct expects a *zerolog.Logger
		us:     storeUser,
		as:     storeArticle,
	}

	tests := []testCase{
		{
			name: "Successful Retrieval of Feed Articles",
			ctx:  context.WithValue(context.Background(), auth.KeyUserID, uint(1)), // Mock authenticated context
			req:  &pb.GetFeedArticlesRequest{Limit: 10, Offset: 0},
			setup: func(us *store.UserStore, as *store.ArticleStore) {
				mock.ExpectQuery("SELECT .* FROM users WHERE id = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"to_user_id"}).AddRow(2))

				mock.ExpectQuery("SELECT .* FROM articles WHERE user_id in \\(\\?\\)").
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at"}).
						AddRow(1, "Title", time.Now(), time.Now()))
			},
			check: func(resp *pb.ArticlesResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, int32(1), resp.GetArticlesCount())
			},
		},
		{
			name: "Unauthenticated User",
			ctx:  context.Background(), // No user ID in context
			req:  &pb.GetFeedArticlesRequest{Limit: 10, Offset: 0},
			setup: func(us *store.UserStore, as *store.ArticleStore) {
				// No specific setup required for unauthenticated
			},
			check: func(resp *pb.ArticlesResponse, err error) {
				assert.Nil(t, resp)
				assert.Equal(t, codes.Unauthenticated, status.Code(err))
			},
		},
		// Add other scenarios...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(storeUser, storeArticle)
			resp, err := handler.GetFeedArticles(tc.ctx, tc.req)
			tc.check(resp, err)
		})
	}
}
