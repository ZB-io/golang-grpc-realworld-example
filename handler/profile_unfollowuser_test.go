package handler

import (
	"context"
	"testing"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/store"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
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
	type test struct {
		name             string
		ctx              context.Context
		req              *pb.UnfollowRequest
		setupMocks       func(us *store.UserStore, mock sqlmock.Sqlmock)
		expectedResponse *pb.ProfileResponse
		expectedError    codes.Code
	}

	tests := []test{
		{
			name: "Successfully Unfollow a User",
			ctx:  context.WithValue(context.Background(), auth.ContextKey("userID"), uint(1)),
			req:  &pb.UnfollowRequest{Username: "targetUser"},
			setupMocks: func(us *store.UserStore, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentUser"))
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username=").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(2, "targetUser"))
				mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectExec("DELETE FROM (.+) WHERE").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedResponse: &pb.ProfileResponse{Profile: &pb.Profile{Username: "targetUser", Following: false}},
			expectedError:    codes.OK,
		},
		{
			name: "User Tries to Unfollow Themselves",
			ctx:  context.WithValue(context.Background(), auth.ContextKey("userID"), uint(1)),
			req:  &pb.UnfollowRequest{Username: "currentUser"},
			setupMocks: func(us *store.UserStore, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentUser"))
			},
			expectedError: codes.InvalidArgument,
		},
		{
			name: "User Tries to Unfollow a Nonexistent User",
			ctx:  context.WithValue(context.Background(), auth.ContextKey("userID"), uint(1)),
			req:  &pb.UnfollowRequest{Username: "nonexistentUser"},
			setupMocks: func(us *store.UserStore, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentUser"))
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username=").WillReturnError(errors.New("record not found"))
			},
			expectedError: codes.NotFound,
		},
		{
			name: "Unauthenticated User Attempts to Unfollow",
			ctx:  context.Background(),
			req:  &pb.UnfollowRequest{Username: "targetUser"},
			setupMocks: func(us *store.UserStore, mock sqlmock.Sqlmock) {

			},
			expectedError: codes.Unauthenticated,
		},
		{
			name: "User Not Following the Target User",
			ctx:  context.WithValue(context.Background(), auth.ContextKey("userID"), uint(1)),
			req:  &pb.UnfollowRequest{Username: "targetUser"},
			setupMocks: func(us *store.UserStore, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentUser"))
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username=").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(2, "targetUser"))
				mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedError: codes.Unauthenticated,
		},
		{
			name: "Internal Error During Unfollowing Process",
			ctx:  context.WithValue(context.Background(), auth.ContextKey("userID"), uint(1)),
			req:  &pb.UnfollowRequest{Username: "targetUser"},
			setupMocks: func(us *store.UserStore, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentUser"))
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username=").WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(2, "targetUser"))
				mock.ExpectQuery("SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectExec("DELETE FROM (.+) WHERE").WillReturnError(errors.New("internal server error"))
			},
			expectedError: codes.Aborted,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to open mock SQL connection: %s", err)
			}
			defer db.Close()

			us := &store.UserStore{DB: db}
			logger := zerolog.New(zerolog.ConsoleWriter{Out: nil})
			h := &Handler{logger: &logger, us: us}

			if tc.setupMocks != nil {
				tc.setupMocks(us, mock)
			}

			response, err := h.UnfollowUser(tc.ctx, tc.req)

			if tc.expectedError == codes.OK {
				if response == nil || response.Profile == nil || response.Profile.Username != tc.req.Username {
					t.Errorf("expected valid ProfileResponse, got %v", response)
				} else {
					t.Logf("test scenario '%s' succeeded", tc.name)
				}
			} else if err == nil || status.Code(err) != tc.expectedError {
				t.Errorf("expected error code %v, got %v", tc.expectedError, status.Code(err))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unmet expectations: %s", err)
			}
		})
	}
}
