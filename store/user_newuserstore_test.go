package store

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"testing"
	"reflect"
)

func TestNewUserStore(t *testing.T) {

	// initialize a list of test cases (table-driven tests)
	cases := []struct {
		name string
		db   *gorm.DB
	}{
		{
			"Creation of NewUserStore with Valid DB",
			// TODO: Replace the nil mock DB with a valid configured one.
			&gorm.DB{},
		},
		{
			"Creation of NewUserStore with Nil DB",
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Executing test scenario:", tc.name)

			db, _, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}

			us := NewUserStore(gdb)

			// validate that the returned UserStore's db matches the DB passed
			if us.db != gdb {
				t.Errorf("Expected DB in UserStore to be %v, but got %v", gdb, us.db)
			}

			if tc.db == nil && us.db != nil {
				t.Errorf("Expected DB in UserStore to be %v, but got %v", tc.db, us.db)
			}

			t.Log("Success: Expected DB matches the DB in UserStore for test scenario:", tc.name)
		})
	}

	// Repeat the creation
	t.Run("Repeated Creation of NewUserStore with Same DB", func(t *testing.T) {
		t.Log("Executing test scenario: Repeated Creation of NewUserStore with Same DB")

		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
		}

		gdb, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
		}

		us1 := NewUserStore(gdb)
		us2 := NewUserStore(gdb)

		// Validate that two UserStores are different instances
		if reflect.DeepEqual(us1, us2) {
			t.Errorf("Expected two different instances of UserStore, but they are the same")
		}

		t.Log("Success: Two different instances of UserStore created for the same DB")
	})
}
