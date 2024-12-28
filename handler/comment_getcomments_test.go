package handler

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlerGetComments(t *testing.T) {
	type testCase struct {
		description  string
		setupMock    func(sqlmock.Sqlmock)
		request      *pb.GetCommentsRequest
		expectedErr  codes.Code
		expectedResp *pb.CommentsResponse
	}

	tests := []testCase{
		{
			description: "Scenario 1: Successfully Retrieve Comments for a Valid Article",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" WHERE`).
					WithArgs(123).WillReturnRows(sqlmock.NewRows([]string{"id", "slug"}).
					AddRow(123, "test-slug"))

				mock.ExpectQuery(`SELECT (.+) FROM "comments" WHERE "article_id" = ?`).
					WithArgs(123).WillReturnRows(sqlmock.NewRows([]string{"id", "body"}).
					AddRow(1, "test comment"))

				mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
					AddRow(1, "user123"))

				mock.ExpectQuery(`SELECT count(.+) FROM "follows" WHERE`).
					WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			request: &pb.GetCommentsRequest{
				Slug: "123",
			},
			expectedErr: codes.OK,
			expectedResp: &pb.CommentsResponse{
				Comments: []*pb.Comment{
					{
						Id:   "1",
						Body: "test comment",
						Author: &pb.Profile{
							Username:  "user123",
							Following: false,
						},
					},
				},
			},
		},
		{
			description: "Scenario 2: Invalid Slug Type (Non-integer Slug)",
			setupMock:   func(mock sqlmock.Sqlmock) {},
			request: &pb.GetCommentsRequest{
				Slug: "invalid-slug",
			},
			expectedErr: codes.InvalidArgument,
			expectedResp: nil,
		},
		{
			description: "Scenario 3: Article Not Found for Given Slug",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" WHERE`).
					WithArgs(999).WillReturnError(sqlmock.ErrNoRows)
			},
			request: &pb.GetCommentsRequest{
				Slug: "999",
			},
			expectedErr: codes.InvalidArgument,
			expectedResp: nil,
		},
		{
			description: "Scenario 4: Error Occurs When Fetching Comments",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" WHERE`).
					WithArgs(124).WillReturnRows(sqlmock.NewRows([]string{"id", "slug"}).
					AddRow(124, "test-slug"))

				mock.ExpectQuery(`SELECT (.+) FROM "comments" WHERE "article_id" = ?`).
					WithArgs(124).WillReturnError(sqlmock.ErrCancelled)
			},
			request: &pb.GetCommentsRequest{
				Slug: "124",
			},
			expectedErr: codes.Aborted,
			expectedResp: nil,
		},
		{
			description: "Scenario 5: User Not Found when Resolving Current User",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" WHERE`).
					WithArgs(125).WillReturnRows(sqlmock.NewRows([]string{"id", "slug"}).
					AddRow(125, "test-slug"))

				mock.ExpectQuery(`SELECT (.+) FROM "comments" WHERE "article_id" = ?`).
					WithArgs(125).WillReturnRows(sqlmock.NewRows([]string{"id", "author_id", "body"}).
					AddRow(1, 1, "test comment"))

				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 2, nil
				}

				mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
					WithArgs(2).WillReturnError(sqlmock.ErrNoRows)
			},
			request: &pb.GetCommentsRequest{
				Slug: "125",
			},
			expectedErr: codes.NotFound,
			expectedResp: nil,
		},
		{
			description: "Scenario 6: Following Status Retrieval Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT (.+) FROM "articles" WHERE`).
					WithArgs(126).WillReturnRows(sqlmock.NewRows([]string{"id", "slug"}).
					AddRow(126, "test-slug"))

				// Assume error when fetching isFollowing status
				mock.ExpectQuery(`SELECT (.+) FROM "comments" WHERE "article_id" = ?`).
					WithArgs(126).WillReturnRows(sqlmock.NewRows([]string{"id", "author_id", "body"}).
					AddRow(1, 1, "test comment"))

				auth.GetUserID = func(ctx context.Context) (uint, error) {
					return 1, nil
				}

				mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
					AddRow(1, "user123"))

				mock.ExpectQuery(`SELECT count(.+) FROM "follows" WHERE`).
					WithArgs(1, 1).WillReturnError(sqlmock.ErrCancelled)
			},
			request: &pb.GetCommentsRequest{
				Slug: "126",
			},
			expectedErr: codes.NotFound,
			expectedResp: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
				return
			}
			defer db.Close()

			tc.setupMock(mock)

			handler := &Handler{
				logger: nil, // Initialize the logger appropriately. TODO: Set up logger for detailed logging
				us:     store.NewUserStore(db),
				as:     store.NewArticleStore(db),
			}

			resp, err := handler.GetComments(context.Background(), tc.request)
			if err != nil {
				if tc.expectedErr == codes.OK {
					t.Fatalf("unexpected error: %v", err)
				}
				if status.Code(err) != tc.expectedErr {
					t.Fatalf("expected error: %v, got: %v", tc.expectedErr, status.Code(err))
				}
			} else {
				if tc.expectedErr != codes.OK {
					t.Fatalf("expected error: %v, got: nil", tc.expectedErr)
				}

				if len(resp.Comments) != len(tc.expectedResp.Comments) {
					t.Fatalf("unexpected number of comments. got: %v want: %v", len(resp.Comments), len(tc.expectedResp.Comments))
				}

				// Further validation can be added to compare the actual response to the expected response in detail
			}
		})
	}
}
