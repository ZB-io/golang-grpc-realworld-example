package store

import (
	"testing"
	"errors"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	// Create a slice of test cases
	testCases := []struct {
		name          string
		user          model.User
		expectedIDs   []uint
		mockDB        func(mock sqlmock.Sqlmock, user model.User)
		expectedError error
	}{
		{
			name: "Valid User with Following Users",
			user: model.User{Model: gorm.Model{ID: 1}},
			expectedIDs: []uint{2, 3, 4},
			mockDB: func(mock sqlmock.Sqlmock, user model.User) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)
				mock.ExpectQuery(`^SELECT to_user_id FROM follows WHERE from_user_id = \?`).
					WithArgs(user.ID).
					WillReturnRows(rows)
			},
			expectedError: nil,
		},
		{
			name: "Valid User with No Following Users",
			user: model.User{Model: gorm.Model{ID: 1}},
			expectedIDs: []uint{},
			mockDB: func(mock sqlmock.Sqlmock, user model.User) {
				mock.ExpectQuery(`^SELECT to_user_id FROM follows WHERE from_user_id = \?`).
					WithArgs(user.ID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: nil,
		},
		{
			name: "Invalid User",
			user: model.User{Model: gorm.Model{ID: 0}},
			expectedIDs: []uint{},
			mockDB: func(mock sqlmock.Sqlmock, user model.User) {
				mock.ExpectQuery(`^SELECT to_user_id FROM follows WHERE from_user_id = \?`).
					WithArgs(user.ID).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Error",
			user: model.User{Model: gorm.Model{ID: 1}},
			expectedIDs: []uint{},
			mockDB: func(mock sqlmock.Sqlmock, user model.User) {
				mock.ExpectQuery(`^SELECT to_user_id FROM follows WHERE from_user_id = \?`).
					WithArgs(user.ID).
					WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the database
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			tc.mockDB(mock, tc.user)

			// Create a new UserStore
			store := &UserStore{db: gormDB}

			// Invoke the function
			ids, err := store.GetFollowingUserIDs(&tc.user)

			// Assert function return
			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error '%s', got '%s'", tc.expectedError, err)
			}

			if len(ids) != len(tc.expectedIDs) {
				t.Errorf("expected length of ids '%d', got '%d'", len(tc.expectedIDs), len(ids))
			}

			for i := range ids {
				if ids[i] != tc.expectedIDs[i] {
					t.Errorf("expected id '%d', got '%d'", tc.expectedIDs[i], ids[i])
				}
			}
		})
	}
}
