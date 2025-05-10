package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
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

type Rows struct {
	converter driver.ValueConverter
	cols      []string
	def       []*Column
	rows      [][]driver.Value
	pos       int
	nextErr   map[int]error
	closeErr  error
}

type Controller struct {
	mu            sync.Mutex
	t             TestReporter
	expectedCalls *callSet
	finished      bool
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
func TestHandlerFollowUser(t *testing.T) {
	tests := []struct {
		name          string
		mockUserID    func() (uint, error)
		mockGetByID   func(m sqlmock.Sqlmock)
		mockGetByUser func(m sqlmock.Sqlmock)
		mockFollow    func(m sqlmock.Sqlmock)
		request       *pb.FollowRequest
		expectedErr   error
	}{
		{
			name: "Scenario 1: User is Unauthenticated",
			mockUserID: func() (uint, error) {
				return 0, errors.New("unauthenticated")
			},
			request:     &pb.FollowRequest{Username: "targetuser"},
			expectedErr: status.Error(codes.Unauthenticated, "unauthenticated"),
		},
		{
			name: "Scenario 2: Current User Not Found",
			mockUserID: func() (uint, error) {
				return 1, nil
			},
			mockGetByID: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT \\* FROM `users` WHERE").WillReturnError(errors.New("not found"))
			},
			request:     &pb.FollowRequest{Username: "targetuser"},
			expectedErr: status.Error(codes.NotFound, "user not found"),
		},
		{
			name: "Scenario 3: User Attempts to Follow Themselves",
			mockUserID: func() (uint, error) {
				return 1, nil
			},
			mockGetByID: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "selfuser")
				m.ExpectQuery("SELECT \\* FROM `users` WHERE").WillReturnRows(rows)
			},
			request:     &pb.FollowRequest{Username: "selfuser"},
			expectedErr: status.Error(codes.InvalidArgument, "cannot follow yourself"),
		},
		{
			name: "Scenario 4: Target User Not Found",
			mockUserID: func() (uint, error) {
				return 1, nil
			},
			mockGetByID: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentuser")
				m.ExpectQuery("SELECT \\* FROM `users` WHERE").WillReturnRows(rows)
			},
			mockGetByUser: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT \\* FROM `users` WHERE `username`=\\?").WillReturnError(errors.New("user not found"))
			},
			request:     &pb.FollowRequest{Username: "targetuser"},
			expectedErr: status.Error(codes.NotFound, "user was not found"),
		},
		{
			name: "Scenario 5: Successful Follow Operation",
			mockUserID: func() (uint, error) {
				return 1, nil
			},
			mockGetByID: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "currentuser")
				m.ExpectQuery("SELECT \\* FROM `users` WHERE").WillReturnRows(rows)
			},
			mockGetByUser: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(2, "targetuser")
				m.ExpectQuery("SELECT \\* FROM `users` WHERE `username`=\\?").WillReturnRows(rows)
			},
			mockFollow: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO `follows`").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			request:     &pb.FollowRequest{Username: "targetuser"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			mockUserStore := store.NewMockUserStore(ctrl)
			mockUserStore.EXPECT().GetByID(gomock.Any()).DoAndReturn(func(id uint) (*model.User, error) {
				if tt.mockGetByID != nil {
					tt.mockGetByID(mock)
				}
				userMock := model.User{ID: id, Username: "mockuser"}
				return &userMock, nil
			}).AnyTimes()

			mockUserStore.EXPECT().GetByUsername(gomock.Any()).DoAndReturn(func(username string) (*model.User, error) {
				if tt.mockGetByUser != nil {
					tt.mockGetByUser(mock)
				}
				userMock := model.User{ID: 1, Username: username}
				return &userMock, nil
			}).AnyTimes()

			mockUserStore.EXPECT().Follow(gomock.Any(), gomock.Any()).DoAndReturn(func(a *model.User, b *model.User) error {
				if tt.mockFollow != nil {
					tt.mockFollow(mock)
				}
				return nil
			}).AnyTimes()

			logger := zerolog.Logger{}
			h := &Handler{
				logger: &logger,
				us:     mockUserStore,
			}

			auth.GetUserID = tt.mockUserID

			resp, err := h.FollowUser(context.Background(), tt.request)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
