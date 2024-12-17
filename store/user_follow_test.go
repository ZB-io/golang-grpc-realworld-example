package store

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/require"
)

func setupTestDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDB, gormErr := gorm.Open("postgres", mockDB)
	if gormErr != nil {
		return nil, nil, gormErr
	}
	return gormDB, mock, nil
}

func TestUserStoreFollow(t *testing.T) {
	// Mock database setup
	gormDB, _, err := setupTestDB()
	if err != nil {
		t.Fatal("failed to setup test database: ", err)
	}
	store := &UserStore{gormDB}

	var scenarios = []struct {
		name     string
		setup    func(a, b *model.User)
		userA    model.User
		userB    model.User
		expected error
	}{
		{
			name: "Happy Path - when User A is able to follow User B",
			setup: func(a, b *model.User) {
				// Nothing to setup for this case
			},
			userA:    model.User{Username: "John Doe", Email: "john@doe.com", Password: "password"},
			userB:    model.User{Username: "Jane Doe", Email: "jane@doe.com", Password: "password"},
			expected: nil,
		},
		{
			name: "Failed follow attempt - Follows list remains unchanged when following a user fails",
			setup: func(a, b *model.User) {
				// Fill a's follows with some dummy data
				a.Follows = append(a.Follows, model.User{Username: "Dummy", Email: "dummy@user.com", Password: "password"})
			},
			userA:    model.User{Username: "John Doe", Email: "john@doe.com", Password: "password"},
			userB:    model.User{Username: "Jane Doe", Email: "jane@doe.com", Password: "password"},
			expected: gorm.ErrRecordNotFound,
		},
		{
			name: "Attempt to follow nonexistent user - Error handling when User isn't in the database",
			setup: func(a, b *model.User) {
				// Nothing to setup for this case
			},
			userA:    model.User{Username: "John Doe", Email: "john@doe.com", Password: "password"},
			userB:    model.User{Username: "Jane Doe", Email: "jane@doe.com", Password: "password"},
			expected: gorm.ErrRecordNotFound,
		},
		{
			name: "Attempt to follow oneself - User cannot follow himself",
			setup: func(a, b *model.User) {
				// Set a and b to be same user
				*b = *a
			},
			userA:    model.User{Username: "John Doe", Email: "john@doe.com", Password: "password"},
			userB:    model.User{Username: "John Doe", Email: "john@doe.com", Password: "password"}, // TODO: Duplicate user
			expected: errors.New("user cannot follow themselves"),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.setup(&scenario.userA, &scenario.userB)
			err := store.Follow(&scenario.userA, &scenario.userB)
			if scenario.expected != nil {
				require.Error(t, err)
				require.Equal(t, scenario.expected, err)
			} else {
				require.NoError(t, err)
				require.Contains(t, scenario.userA.Follows, scenario.userB)
			}
		})
	}
}
