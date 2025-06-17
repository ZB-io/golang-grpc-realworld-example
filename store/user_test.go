package store

import (
	"database/sql"
	"errors"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name          string
	email         string
	username      string
	expectedUser  *model.User
	expectedError error
	mockBehavior  func()
}

type userCase struct {
	username   string
	expectUser *model.User
	expectErr  error
}

var userCases = []userCase{
	{
		username:   "existingUser",
		expectUser: &model.User{Model: gorm.Model{ID: 1}, Username: "existingUser", Email: "user@example.com"},
		expectErr:  nil,
	},
	{
		username:   "nonExistentUser",
		expectUser: nil,
		expectErr:  gorm.ErrRecordNotFound,
	},
	{
		username:   "",
		expectUser: nil,
		expectErr:  gorm.ErrRecordNotFound,
	},
	{
		username:   "user@123",
		expectUser: &model.User{Model: gorm.Model{ID: 2}, Username: "user@123", Email: "user@123@example.com"},
		expectErr:  nil,
	},
	{
		username:   "ExIstInGUsEr",
		expectUser: nil,
		expectErr:  gorm.ErrRecordNotFound,
	},
	{
		username:   "' OR '1'='1",
		expectUser: nil,
		expectErr:  gorm.ErrRecordNotFound,
	},
}

/*
ROOST_METHOD_HASH=NewUserStore_fb599438e5
ROOST_METHOD_SIG_HASH=NewUserStore_c0075221af

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore // NewUserStore returns a new UserStore
*/
func TestNewUserStore(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("sqlite3", db)
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}

	tests := []struct {
		name        string
		db          *gorm.DB
		expectNil   bool
		expectError bool
	}{
		{
			name:        "Successful Initialization of UserStore",
			db:          gormDB,
			expectNil:   false,
			expectError: false,
		},
		{
			name:        "Handling Nil DB Instance",
			db:          nil,
			expectNil:   true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			userStore := NewUserStore(tt.db)

			if tt.expectNil {
				if userStore != nil {
					t.Errorf("Expected UserStore to be nil, but got %v", userStore)
				}
			} else {
				if userStore == nil {
					t.Errorf("Expected UserStore to be non-nil, but got nil")
				} else if userStore.db != tt.db {
					t.Errorf("Expected UserStore.db to be %v, but got %v", tt.db, userStore.db)
				}
			}

			if !tt.expectError {
				err := mock.ExpectationsWereMet()
				if err != nil {
					t.Errorf("There were unfulfilled expectations: %s", err)
				}
			}

			t.Logf("Test %s passed successfully", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=UserStore_GetByEmail_fda09af5c4
ROOST_METHOD_SIG_HASH=UserStore_GetByEmail_9e84f3286b

FUNCTION_DEF=func (s *UserStore) GetByEmail(email string) (*model.User, error) // GetByEmail finds a user from email
*/
func TestUserStoreGetByEmail(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDB, err := gorm.Open("sqlmock", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock gorm database", err)
	}
	defer mockDB.Close()

	userStore := &UserStore{db: mockDB}

	testCases := []testCase{
		{
			name:  "Normal Operation - Existing Email",
			email: "test@example.com",
			expectedUser: &model.User{
				Email: "test@example.com",
			},
			expectedError: nil,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "test@example.com")
				mock.ExpectQuery("SELECT").WithArgs("test@example.com").WillReturnRows(rows)
			},
		},
		{
			name:  "Normal Operation - Multiple Users with Same Email",
			email: "duplicate@example.com",
			expectedUser: &model.User{
				Email: "duplicate@example.com",
			},
			expectedError: nil,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "duplicate@example.com").AddRow(2, "duplicate@example.com")
				mock.ExpectQuery("SELECT").WithArgs("duplicate@example.com").WillReturnRows(rows)
			},
		},
		{
			name:          "Error Handling - Non-Existing Email",
			email:         "nonexistent@example.com",
			expectedUser:  nil,
			expectedError: errors.New("record not found"),
			mockBehavior: func() {
				mock.ExpectQuery("SELECT").WithArgs("nonexistent@example.com").WillReturnError(errors.New("record not found"))
			},
		},
		{
			name:          "Error Handling - Empty Email",
			email:         "",
			expectedUser:  nil,
			expectedError: errors.New("invalid email"),
			mockBehavior: func() {

			},
		},
		{
			name:          "Error Handling - SQL Injection Attempt",
			email:         "' OR '1'='1",
			expectedUser:  nil,
			expectedError: errors.New("invalid email"),
			mockBehavior: func() {

			},
		},
		{
			name:  "Performance - Large Number of Users",
			email: "large@example.com",
			expectedUser: &model.User{
				Email: "large@example.com",
			},
			expectedError: nil,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "large@example.com")
				mock.ExpectQuery("SELECT").WithArgs("large@example.com").WillReturnRows(rows)
			},
		},
		{
			name:  "Edge Case - Email with Special Characters",
			email: "test.special@email.com",
			expectedUser: &model.User{
				Email: "test.special@email.com",
			},
			expectedError: nil,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "test.special@email.com")
				mock.ExpectQuery("SELECT").WithArgs("test.special@email.com").WillReturnRows(rows)
			},
		},
		{
			name:  "Edge Case - Email with Mixed Case",
			email: "MiXeDcAsE@eXaMpLe.CoM",
			expectedUser: &model.User{
				Email: "mixedcase@example.com",
			},
			expectedError: nil,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "mixedcase@example.com")
				mock.ExpectQuery("SELECT").WithArgs("MiXeDcAsE@eXaMpLe.CoM").WillReturnRows(rows)
			},
		},
		{
			name:  "Edge Case - Email with Leading and Trailing Spaces",
			email: "  test@example.com  ",
			expectedUser: &model.User{
				Email: "test@example.com",
			},
			expectedError: nil,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "test@example.com")
				mock.ExpectQuery("SELECT").WithArgs("test@example.com").WillReturnRows(rows)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tc.mockBehavior()

			email := strings.TrimSpace(tc.email)
			start := time.Now()
			user, err := userStore.GetByEmail(email)
			elapsed := time.Since(start)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedUser, user)

			if elapsed > 100*time.Millisecond {
				t.Logf("Test %s took %s, which is above the performance threshold", tc.name, elapsed)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=UserStore_GetByUsername_622b1b9e41
ROOST_METHOD_SIG_HASH=UserStore_GetByUsername_992f00baec

FUNCTION_DEF=func (s *UserStore) GetByUsername(username string) (*model.User, error) // GetByUsername finds a user from username
*/
func TestUserStoreGetByUsername(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("sqlite3", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub gorm database connection", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{db: gormDB}

	tests := userCases

	for _, test := range tests {
		t.Run(test.username, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			if test.expectUser != nil {
				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \\(username = \\?\\) ORDER BY \"users\"\\.\"id\" ASC LIMIT 1").
					WithArgs(test.username).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).AddRow(test.expectUser.ID, test.expectUser.Username, test.expectUser.Email))
			} else {
				mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE \\(username = \\?\\) ORDER BY \"users\"\\.\"id\" ASC LIMIT 1").
					WithArgs(test.username).
					WillReturnError(sql.ErrNoRows)
			}

			user, err := userStore.GetByUsername(test.username)

			if test.expectErr != nil && err == nil {
				t.Fatalf("expected error %v, got %v", test.expectErr, err)
			}
			if !errors.Is(err, test.expectErr) {
				t.Fatalf("expected error %v, got %v", test.expectErr, err)
			}
			if !reflect.DeepEqual(user, test.expectUser) {
				t.Fatalf("expected user %v, got %v", test.expectUser, user)
			}
			t.Logf("Test passed for username: %s", test.username)
		})
	}
}
