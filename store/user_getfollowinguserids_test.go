package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

func TestUserStoreGetFollowingUserIDs(t *testing.T) {

	type testScenario struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		user          *model.User
		expectedIDs   []uint
		expectedError error
	}

	scenarios := []testScenario{
		{
			name: "User with One Followed User",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"to_user_id"}).AddRow(2))
			},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedIDs:   []uint{2},
			expectedError: nil,
		},
		{
			name: "User with No Followed Users",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"to_user_id"}))
			},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name: "User Following Multiple Users",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnRows(
						sqlmock.NewRows([]string{"to_user_id"}).
							AddRow(2).AddRow(3).AddRow(4),
					)
			},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: nil,
		},
		{
			name: "Non-existing User Input",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(9999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			user:          &model.User{Model: gorm.Model{ID: 9999}},
			expectedIDs:   []uint{},
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Error Occurrence",

			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedIDs:   []uint{},
			expectedError: errors.New("database error"),
		},
		{
			name: "Mixed State of Following",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"to_user_id"}).AddRow(3).AddRow(5))
			},
			user:          &model.User{Model: gorm.Model{ID: 2}},
			expectedIDs:   []uint{3, 5},
			expectedError: nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gdb, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("failed to open gorm db: %v", err)
			}

			scenario.setupMock(mock)

			store := &UserStore{db: gdb}

			ids, err := store.GetFollowingUserIDs(scenario.user)

			t.Logf("Scenario: %s", scenario.name)

			if scenario.expectedError != nil {
				if err == nil || err.Error() != scenario.expectedError.Error() {
					t.Errorf("expected error %v, got %v", scenario.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if !equal(ids, scenario.expectedIDs) {
				t.Errorf("expected IDs %v, got %v", scenario.expectedIDs, ids)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}
func equal(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
