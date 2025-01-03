package handler

import (
	"context"
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

// Mock getUserID function used in auth package
var mockGetUserID = func(ctx context.Context) (uint, error) {
	// TODO: Update this mock with realistic test case scenarios, e.g., return IDs based on context keys
	return 1, nil
}

func TestCreateComment(t *testing.T) {
	// Initialize a new zerolog logger instance
	logger := zerolog.New(zerolog.ConsoleWriter{Out: t}).With().Logger()

	// Mock the database using sqlmock
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening a stub database connection: %s", err)
	}
	defer db.Close()

	// Create mock stores
	userStore := &store.UserStore{Db: db}
	articleStore := &store.ArticleStore{Db: db}

	// Create handler with mock logger and stores
	h := &Handler{logger: &logger, us: userStore, as: articleStore}

	tests := []struct {
		name           string
		ctx            context.Context
		request        *pb.CreateCommentRequest
		expectedError  error
		expectedResult *pb.CommentResponse
	}{
		{
			name: "Successful Comment Creation",
			ctx:  context.TODO(), // Use a context with proper authentication details in real test
			request: &pb.CreateCommentRequest{
				Slug: "1", // Assume this correctly maps to an article
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "This is a test comment",
				},
			},
			expectedError: nil,
			expectedResult: &pb.CommentResponse{
				Comment: &pb.Comment{
					Body: "This is a test comment",
					// Fill additional fields based on realistic response data
				},
			},
		},
		{
			name:           "Unauthenticated User",
			ctx:            context.TODO(), // No authentication details included
			request:        &pb.CreateCommentRequest{}, // Any req object
			expectedError:  status.Error(codes.Unauthenticated, "unauthenticated"),
			expectedResult: nil,
		},
		{
			name: "User Not Found",
			ctx:  context.TODO(), // Mock context with potential proper authentication
			request: &pb.CreateCommentRequest{
				Slug: "1",
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "Comment body",
				},
			},
			expectedError:  status.Error(codes.NotFound, "user not found"),
			expectedResult: nil,
		},
		// Additional test cases would follow the structure above
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the GetUserID to return user ID based on scenario
			auth.GetUserID = mockGetUserID

			response, err := h.CreateComment(tt.ctx, tt.request)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, response)
		})
	}
}
