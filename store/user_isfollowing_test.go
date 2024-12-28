package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"

	"github.com/stretchr/testify/assert"
)

func TestUserStoreIsFollowing(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		setupMocks  func(sqlmock.Sqlmock)
		expectedErr error
		expectedRes bool
	}{
		{
			name: "Scenario 1: Successfully determine when User A is following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}}, // User A with ID 1
			userB: &model.User{Model: gorm.Model{ID: 2}}, // User B with ID 2
			setupMocks: func(mock sqlmock.Sqlmock) {
				// Mock database response to simulate follow relationship
				mock.ExpectQuery(`SELECT count\(.+\) FROM "follows"`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedErr: nil,
			expectedRes: true,
		},
		{
			name: "Scenario 2: Successfully determine when User A is not following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 3}},
			setupMocks: func(mock sqlmock.Sqlmock) {
				// Mock DB response for no follow relationship
				mock.ExpectQuery(`SELECT count\(.+\) FROM "follows"`).
					WithArgs(1, 3).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedErr: nil,
			expectedRes: false,
		},
		{
			name: "Scenario 3: Handle nil User A parameter",
			userA: nil,
			userB: &model.User{Model: gorm.Model{ID: 1}},
			setupMocks: func(mock sqlmock.Sqlmock) {
				// No DB call expected since userA is nil
			},
			expectedErr: nil,
			expectedRes: false,
		},
		{
			name: "Scenario 4: Handle nil User B parameter",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: nil,
			setupMocks: func(mock sqlmock.Sqlmock) {
				// No DB call expected since userB is nil
			},
			expectedErr: nil,
			expectedRes: false,
		},
		{
			name: "Scenario 5: SQL error while querying the database",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMocks: func(mock sqlmock.Sqlmock) {
				// Simulate SQL error
				mock.ExpectQuery(`SELECT count\(.+\) FROM "follows"`).
					WithArgs(1, 2).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedErr: gorm.ErrInvalidSQL,
			expectedRes: false,
		},
		{
			name: "Scenario 6: Users exist but the 'follows' table is empty",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMocks: func(mock sqlmock.Sqlmock) {
				// No relationship, empty follows table
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
			// Create a mock DB
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			assert.NoError(t, err)

			userStore := &UserStore{db: gormDB}

			// Set up mock expectations
			tt.setupMocks(mock)

			// Run the test case
			res, err := userStore.IsFollowing(tt.userA, tt.userB)

			// Make assertions
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedRes, res)

			// Check if all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)

			t.Logf("Test case '%s' passed", tt.name)
		})
	}
}
