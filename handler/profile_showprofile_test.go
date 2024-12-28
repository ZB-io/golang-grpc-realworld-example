package handler

import (
	"context"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
func TestHandlerShowProfile(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name            string
		contextSetup    func() context.Context
		requestSetup    func() *pb.ShowProfileRequest
		userStoreSetup  func(mock sqlmock.Sqlmock)
		expectError     bool
		expectedCode    codes.Code
		expectedProfile *pb.Profile
	}

	tests := []testCase{
		{
			name: "Scenario 1: Successfully Show Profile of a User",
			contextSetup: func() context.Context {

				return auth.ContextWithUserID(context.Background(), 1)
			},
			requestSetup: func() *pb.ShowProfileRequest {
				return &pb.ShowProfileRequest{Username: "validUser"}
			},
			userStoreSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentUser"))

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."username" = \$1`).
					WithArgs("validUser").
					WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("validUser"))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "follows"`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectError:  false,
			expectedCode: codes.OK,
			expectedProfile: &pb.Profile{
				Username:  "validUser",
				Following: true,
			},
		},
		{
			name: "Scenario 2: Unauthenticated User Trying to Show Profile",
			contextSetup: func() context.Context {

				return context.Background()
			},
			requestSetup: func() *pb.ShowProfileRequest {
				return &pb.ShowProfileRequest{Username: "validUser"}
			},
			userStoreSetup: func(mock sqlmock.Sqlmock) {},
			expectError:    true,
			expectedCode:   codes.Unauthenticated,
		},
		{
			name: "Scenario 3: Requested User Not Found",
			contextSetup: func() context.Context {
				return auth.ContextWithUserID(context.Background(), 1)
			},
			requestSetup: func() *pb.ShowProfileRequest {
				return &pb.ShowProfileRequest{Username: "nonExistentUser"}
			},
			userStoreSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentUser"))

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."username" = \$1`).
					WithArgs("nonExistentUser").
					WillReturnError(status.Error(codes.NotFound, "user was not found"))
			},
			expectError:  true,
			expectedCode: codes.NotFound,
		},
		{
			name: "Scenario 4: Current User Not Found",
			contextSetup: func() context.Context {
				return auth.ContextWithUserID(context.Background(), 9999)
			},
			requestSetup: func() *pb.ShowProfileRequest {
				return &pb.ShowProfileRequest{Username: "validUser"}
			},
			userStoreSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1`).
					WithArgs(9999).
					WillReturnError(status.Error(codes.NotFound, "user not found"))
			},
			expectError:  true,
			expectedCode: codes.NotFound,
		},
		{
			name: "Scenario 5: Internal Error on Checking Following Status",
			contextSetup: func() context.Context {
				return auth.ContextWithUserID(context.Background(), 1)
			},
			requestSetup: func() *pb.ShowProfileRequest {
				return &pb.ShowProfileRequest{Username: "validUser"}
			},
			userStoreSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentUser"))
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."username" = \$1`).
					WithArgs("validUser").
					WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("validUser"))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "follows"`).
					WillReturnError(status.Error(codes.Internal, "internal server error"))
			},
			expectError:  true,
			expectedCode: codes.NotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			userStore := &store.UserStore{DB: db}
			logger := zerolog.New(os.Stdout)

			handler := &Handler{
				logger: &logger,
				us:     userStore,
			}

			ctx := tc.contextSetup()
			req := tc.requestSetup()
			tc.userStoreSetup(mock)

			resp, err := handler.ShowProfile(ctx, req)
			if tc.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if status.Code(err) != tc.expectedCode {
					t.Fatalf("expected error code %v but got %v", tc.expectedCode, status.Code(err))
				}
				t.Logf("test: %s passed with expected error: %v", tc.name, err)
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if resp.Profile.Username != tc.expectedProfile.Username || resp.Profile.Following != tc.expectedProfile.Following {
					t.Fatalf("unexpected profile data: got %+v", resp.Profile)
				}
				t.Logf("test: %s passed with expected profile: %+v", tc.name, resp.Profile)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
