package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestUserStoreIsFollowing(t *testing.T) {

	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		setupMocks  func(sqlmock.Sqlmock)
		expectedErr error
		expectedRes bool
	}{
		{
			name:  "Scenario 1: Successfully determine when User A is following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMocks: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT count\(.+\) FROM "follows"`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedErr: nil,
			expectedRes: true,
		},
		{
			name:  "Scenario 2: Successfully determine when User A is not following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 3}},
			setupMocks: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT count\(.+\) FROM "follows"`).
					WithArgs(1, 3).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedErr: nil,
			expectedRes: false,
		},
		{
			name:  "Scenario 3: Handle nil User A parameter",
			userA: nil,
			userB: &model.User{Model: gorm.Model{ID: 1}},
			setupMocks: func(mock sqlmock.Sqlmock) {

			},
			expectedErr: nil,
			expectedRes: false,
		},
		{
			name:  "Scenario 4: Handle nil User B parameter",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: nil,
			setupMocks: func(mock sqlmock.Sqlmock) {

			},
			expectedErr: nil,
			expectedRes: false,
		},
		{
			name:  "Scenario 5: SQL error while querying the database",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMocks: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT count\(.+\) FROM "follows"`).
					WithArgs(1, 2).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedErr: gorm.ErrInvalidSQL,
			expectedRes: false,
		},
		{
			name:  "Scenario 6: Users exist but the 'follows' table is empty",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMocks: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT count\(.+\) FROM "follows"`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedErr: nil,
			expectedRes: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			assert.NoError(t, err)

			userStore := &UserStore{db: gormDB}

			tt.setupMocks(mock)

			res, err := userStore.IsFollowing(tt.userA, tt.userB)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedRes, res)

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)

			t.Logf("Test case '%s' passed", tt.name)
		})
	}
}
