package store

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestNewUserStore(t *testing.T) {
	t.Parallel()

	// Define test cases
	testCases := []struct {
		name            string
		db              *gorm.DB
		expectedDBEqual bool
	}{
		{
			name:            "Successful creation of NewUserStore",
			db:              new(gorm.DB),
			expectedDBEqual: true,
		},
		{
			name:            "Creation of NewUserStore with nil gorm.DB instance",
			db:              nil,
			expectedDBEqual: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			userStore := NewUserStore(tc.db)

			// Assert
			if userStore == nil {
				t.Fatal("Expected NewUserStore to return a non-nil UserStore instance, but got nil")
			}

			if (userStore.db == nil) != (!tc.expectedDBEqual) {
				t.Fatalf("Expected UserStore.db to be %v, but got %v", tc.expectedDBEqual, userStore.db != nil)
			}

			if userStore.db != nil && !reflect.DeepEqual(userStore.db, tc.db) {
				t.Fatal("Expected UserStore.db to be equal to the provided gorm.DB instance, but they were different")
			}
		})
	}

	// Scenario 3: Multiple invocations of NewUserStore
	t.Run("Multiple invocations of NewUserStore", func(t *testing.T) {
		// Arrange
		db, _, _ := sqlmock.New()
		gormDB, _ := gorm.Open("postgres", db)

		// Act
		userStore1 := NewUserStore(gormDB)
		userStore2 := NewUserStore(gormDB)

		// Assert
		if userStore1 == userStore2 {
			t.Fatal("Expected multiple invocations of NewUserStore to return distinct UserStore instances, but they were the same")
		}

		if userStore1.db != userStore2.db {
			t.Fatal("Expected UserStore.db of both instances to be equal to the provided gorm.DB instance, but they were different")
		}
	})
}
