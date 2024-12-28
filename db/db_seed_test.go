package db

import (
	"errors"
	"fmt"
	"testing"
	"github.com/BurntSushi/toml"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"io/ioutil"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestSeed(t *testing.T) {
	scenarios := []struct {
		Name          string
		TomlContent   string
		ExpectError   bool
		ExpectedCount int
		SetupMock     func(mock sqlmock.Sqlmock)
	}{
		{
			Name: "Successfully Seed Users into the Database",
			TomlContent: `
				[[Users]]
				Id = 1
				Email = "test1@example.com"

				[[Users]]
				Id = 2
				Email = "test2@example.com"
			`,
			ExpectError:   false,
			ExpectedCount: 2,
			SetupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(2, 1))
			},
		},
		{
			Name: "Fail to Seed Users When the TOML File is Missing",
			TomlContent: `
			`,
			ExpectError:   true,
			ExpectedCount: 0,
			SetupMock:     func(mock sqlmock.Sqlmock) {},
		},
		{
			Name:        "Erroneous TOML Format",
			TomlContent: `[[Users]] Id =`,
			ExpectError: true,
			SetupMock:   func(mock sqlmock.Sqlmock) {},
		},
		{
			Name: "Database Insertion Error",
			TomlContent: `
				[[Users]]
				Id = 1
				Email = "test@example.com"
			`,
			ExpectError:   true,
			ExpectedCount: 0,
			SetupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("db error"))
			},
		},
		{
			Name: "Successful Seeding with No Users in TOML",
			TomlContent: `
				[[Users]]
			`,
			ExpectError:   false,
			ExpectedCount: 0,
			SetupMock:     func(mock sqlmock.Sqlmock) {},
		},
		{
			Name: "Duplicate User Entries in TOML File",
			TomlContent: `
				[[Users]]
				Id = 1
				Email = "test@example.com"

				[[Users]]
				Id = 1
				Email = "test@example.com"
			`,
			ExpectError:   false,
			ExpectedCount: 2,
			SetupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(2, 1))
			},
		},
		{
			Name: "Check Side Effects on DB after Partial Failure",
			TomlContent: `
				[[Users]]
				Id = 1
				Email = "first@example.com"

				[[Users]]
				Id = 2
				Email = "second@example.com"
			`,
			ExpectError:   true,
			ExpectedCount: 1,
			SetupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO").WillReturnError(errors.New("db error"))
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gdb, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("Failed to open gorm db connection: %v", err)
			}

			s.SetupMock(mock)

			ioutilReadFileOriginal := ioutil.ReadFile
			ioutil.ReadFile = func(filename string) ([]byte, error) {
				if s.TomlContent == "" {
					return nil, errors.New("file not found")
				}
				return []byte(s.TomlContent), nil
			}
			defer func() { ioutil.ReadFile = ioutilReadFileOriginal }()

			err = Seed(gdb)

			if s.ExpectError && err == nil {
				t.Errorf("Expected an error but got none.")
			} else if !s.ExpectError && err != nil {
				t.Errorf("Did not expect an error but got one: %v", err)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
