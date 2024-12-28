package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	mockAuth "github.com/raahii/golang-grpc-realworld-example/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func TestHandlerUnfollowUser(t *testing.T) {
	logger := zerolog.New(nil)
	db, mock, _ := sqlmock.New()
	defer db.Close()

	us := &store.UserStore{db: db}
	handler := &Handler{logger: &logger, us: us}

	type testCase struct {
		name       string
		setup      func(mock sqlmock.Sqlmock)
		request    *pb.UnfollowRequest
		expectErr  codes.Code
		expectResp *pb.ProfileResponse
	}

	testCases := []testCase{
		{
			name: "Successfully Unfollow a User",
			setup: func(mock sqlmock.Sqlmock) {
				ctx := context.WithValue(context.Background(), mockAuth.KeyUserID, uint(1))
				currentUser := store.User{ID: 1, Username: "current_user"}
				targetUser := store.User{ID: 2, Username: "user_to_unfollow"}

				mock.ExpectQuery("SELECT").WithArgs(currentUser.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "username"}).AddRow(currentUser.ID, currentUser.Username))

				mock.ExpectQuery("SELECT").WithArgs(targetUser.Username).WillReturnRows(
					sqlmock.NewRows([]string{"id", "username"}).AddRow(targetUser.ID, targetUser.Username))

				mock.ExpectQuery("SELECT count").WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1))

				mock.ExpectExec("DELETE FROM follows").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			request:    &pb.UnfollowRequest{Username: "user_to_unfollow"},
			expectErr:  codes.OK,
			expectResp: &pb.ProfileResponse{Profile: &pb.Profile{Username: "user_to_unfollow", Following: false}},
		},
		{
			name: "User Tries to Unfollow Themselves",
			setup: func(mock sqlmock.Sqlmock) {
				ctx := context.WithValue(context.Background(), mockAuth.KeyUserID, uint(1))
				user := store.User{ID: 1, Username: "current_user"}

				mock.ExpectQuery("SELECT").WithArgs(user.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "username"}).AddRow(user.ID, user.Username))
			},
			request:    &pb.UnfollowRequest{Username: "current_user"},
			expectErr:  codes.InvalidArgument,
			expectResp: nil,
		},
		{
			name: "User Tries to Unfollow a Nonexistent User",
			setup: func(mock sqlmock.Sqlmock) {
				ctx := context.WithValue(context.Background(), mockAuth.KeyUserID, uint(1))
				user := store.User{ID: 1, Username: "current_user"}

				mock.ExpectQuery("SELECT").WithArgs(user.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "username"}).AddRow(user.ID, user.Username))

				mock.ExpectQuery("SELECT").WithArgs("nonexistent_user").WillReturnError(errors.New("no rows found"))
			},
			request:    &pb.UnfollowRequest{Username: "nonexistent_user"},
			expectErr:  codes.NotFound,
			expectResp: nil,
		},
		{
			name: "Unauthenticated User Attempts to Unfollow",
			setup: func(mock sqlmock.Sqlmock) {

			},
			request:    &pb.UnfollowRequest{Username: "any_user"},
			expectErr:  codes.Unauthenticated,
			expectResp: nil,
		},
		{
			name: "User Not Following the Target User",
			setup: func(mock sqlmock.Sqlmock) {
				ctx := context.WithValue(context.Background(), mockAuth.KeyUserID, uint(1))
				currentUser := store.User{ID: 1, Username: "current_user"}
				targetUser := store.User{ID: 2, Username: "not_followed_user"}

				mock.ExpectQuery("SELECT").WithArgs(currentUser.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "username"}).AddRow(currentUser.ID, currentUser.Username))

				mock.ExpectQuery("SELECT").WithArgs(targetUser.Username).WillReturnRows(
					sqlmock.NewRows([]string{"id", "username"}).AddRow(targetUser.ID, targetUser.Username))

				mock.ExpectQuery("SELECT count").WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			request:    &pb.UnfollowRequest{Username: "not_followed_user"},
			expectErr:  codes.Unauthenticated,
			expectResp: nil,
		},
		{
			name: "Internal Error During Unfollowing Process",
			setup: func(mock sqlmock.Sqlmock) {
				ctx := context.WithValue(context.Background(), mockAuth.KeyUserID, uint(1))
				currentUser := store.User{ID: 1, Username: "current_user"}
				targetUser := store.User{ID: 2, Username: "user_to_unfollow"}

				mock.ExpectQuery("SELECT").WithArgs(currentUser.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "username"}).AddRow(currentUser.ID, currentUser.Username))

				mock.ExpectQuery("SELECT").WithArgs(targetUser.Username).WillReturnRows(
					sqlmock.NewRows([]string{"id", "username"}).AddRow(targetUser.ID, targetUser.Username))

				mock.ExpectQuery("SELECT count").WillReturnRows(
					sqlmock.NewRows([]string{"count"}).AddRow(1))

				mock.ExpectExec("DELETE FROM follows").WillReturnError(errors.New("internal error"))
			},
			request:    &pb.UnfollowRequest{Username: "user_to_unfollow"},
			expectErr:  codes.Aborted,
			expectResp: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(mock)

			ctx := context.Background()
			resp, err := handler.UnfollowUser(ctx, tc.request)

			if tc.expectErr != codes.OK {
				assert.Error(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.expectErr, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResp, resp)
			}
		})
	}
}
