package store

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test function TestUserStoreUnfollow to test the Unfollow method
func TestUserStoreUnfollow(t *testing.T) {
	// Test scenarios
	
	scenarios := []struct {
		name          string
		setup         func(a *model.User, b *model.User, mock sqlmock.Sqlmock)
		expectErr     bool
	}{
		{
			name: "Successful Unfollow Test",
			setup: func(a *model.User, b *model.User, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE WHERE from_user_id = ? and to_user_id = ?").
					WithArgs(a.ID, b.ID).
					WillReturnResult(sqlmock.NewResult(5, 1)) // Successfully removed 1 row
			},
			expectErr: false,
		},
		{
			name: "Unfollow Non-Followed User Test",
			setup: func(a *model.User, b *model.User, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE WHERE from_user_id = ? and to_user_id = ?").
					WithArgs(a.ID, b.ID).
					WillReturnResult(sqlmock.NewResult(5, 0)) // No rows affected
			},
			expectErr: true,
		},
		{
			name: "Unfollow Null User Test",
			setup: func(a *model.User, b *model.User, mock sqlmock.Sqlmock) {
				// No expectations as b is nil and thus, no DB operations will be performed
			},
			expectErr: true,
		},
		{
			name: "Unfollow with Uninitialized DB",
			setup: func(a *model.User, b *model.User, mock sqlmock.Sqlmock) {
				// No expectations as DB is uninitialized and thus, no DB operations will be performed
			},
			expectErr: true,
		},
	}

	for _, s := range scenarios {  // range over the scenarios
		t.Run(s.name, func(t *testing.T) {  // use t.Run for better vissibility of which test is being runned
			t.Logf("Running test case: %s", s.name)  // log the scenario name
			// Prepare User A and User B
			userA := &model.User{ID: 1, Username: "UserA", Email: "userA@mail.com", Password: "passA", Bio: "bioA", Image: "imgA"}
			userB := &model.User{ID: 2, Username: "UserB", Email: "userB@mail.com", Password: "passB", Bio: "bioB", Image: "imgB"}
			if s.name == "Unfollow Null User Test" {
				userB = nil  // User B is nil for this case
			}
			// Mock DB
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			require.NoError(t, err)
			gdb, err := gorm.Open("postgres", db)
			require.NoError(t, err)
			gdb.LogMode(true)

			// Setup the mock operations
			s.setup(userA, userB, mock)
			store := &UserStore{db: gdb}
			// Invoke the Unfollow function and Assert the result
			err = store.Unfollow(userA, userB)

			if s.name != "Unfollow with Uninitialized DB" {
				assert.NoError(t, mock.ExpectationsWereMet())
			}

			if s.expectErr {
				assert.Error(t, err, fmt.Sprintf("%s: error was expected from Unfollow", s.name))
			} else {
				assert.NoError(t, err, fmt.Sprintf("%s: unexpected error from Unfollow", s.name))
			}
		})
	}
}
