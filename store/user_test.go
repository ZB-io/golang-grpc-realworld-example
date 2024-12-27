package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920
*/
func TestCreate(t *testing.T) {
	t.Run("Scenario 1: Successfully Create a New User", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"users\"").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		assert.NoError(t, err)

		s := &UserStore{db: gormDB}
		user := &model.User{Username: "testuser", Email: "test@example.com"}

		err = s.Create(user)
		assert.NoError(t, err)
		t.Log("Successfully created user with valid input, no database error.")
	})

	t.Run("Scenario 2: Fail to Create a User Due to Database Error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO \"users\"").
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		assert.NoError(t, err)

		s := &UserStore{db: gormDB}
		user := &model.User{Username: "testuser", Email: "test@example.com"}

		err = s.Create(user)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		t.Logf("Expected database error occurred: %v", err)
	})

	t.Run("Scenario 3: Fail to Create a User with Invalid Data", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		assert.NoError(t, err)

		s := &UserStore{db: gormDB}
		user := &model.User{Username: "", Email: ""}

		err = s.Create(user)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "username or email is empty")
		t.Logf("Invalid data should result in an error: %v", err)
	})

	t.Run("Scenario 4: Handle Nil User Input Gracefully", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		assert.NoError(t, err)

		s := &UserStore{db: gormDB}

		err = s.Create(nil)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "user input is nil")
		t.Log("Function should gracefully handle nil user input and return an error.")
	})
}

/*
ROOST_METHOD_HASH=Follow_48fdf1257b
ROOST_METHOD_SIG_HASH=Follow_8217e61c06
*/
func TestUserStoreFollow(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	userStore := &UserStore{db: gormDB}

	testCases := []struct {
		name        string
		a           *model.User
		b           *model.User
		mockSetup   func()
		expectedErr error
	}{
		{
			name: "Successful Follow Operation",
			a:    &model.User{ID: 1},
			b:    &model.User{ID: 2},
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO follows`).WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Follow Operation with User already Followed",
			a:    &model.User{ID: 1},
			b:    &model.User{ID: 2},
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO follows`).WithArgs(1, 2).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedErr: nil,
		},
		{
			name: "Follow Operation with Database Error",
			a:    &model.User{ID: 1},
			b:    &model.User{ID: 2},
			mockSetup: func() {
				mock.ExpectExec(`INSERT INTO follows`).WithArgs(1, 2).WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedErr: gorm.ErrInvalidSQL,
		},
		{
			name:      "Follow Operation with Null User a",
			a:         nil,
			b:         &model.User{ID: 2},
			mockSetup: func() {},
			expectedErr: errors.New("invalid input"), // Adjust the error type as needed in the implementation
		},
		{
			name:      "Follow Operation with Null User b",
			a:         &model.User{ID: 1},
			b:         nil,
			mockSetup: func() {},
			expectedErr: errors.New("invalid input"), // Adjust the error type as needed in the implementation
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err := userStore.Follow(tc.a, tc.b)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1
*/
func TestUserStoreGetByEmail(t *testing.T) {
	type TestData struct {
		email         string
		expectedUser  *model.User
		expectedError error
		setupMocks    func(sqlmock.Sqlmock)
		scenario      string
	}

	userID := 1
	user := model.User{
		ID:    uint(userID),
		Email: "test@example.com",
		Name:  "Test User",
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	s := &UserStore{db: gormDB}

	tests := []TestData{
		{
			email:         "test@example.com",
			expectedUser:  &user,
			expectedError: nil,
			setupMocks: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "name"}).
					AddRow(user.ID, user.Email, user.Name)
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs(user.Email).
					WillReturnRows(rows)
			},
			scenario: "Successfully Retrieve a User by Email",
		},
		{
			email:         "notfound@example.com",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs("notfound@example.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			scenario: "User Not Found",
		},
		{
			email:         "test@example.com",
			expectedUser:  nil,
			expectedError: errors.New("db connection error"),
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs(user.Email).
					WillReturnError(errors.New("db connection error"))
			},
			scenario: "Database Connection Error",
		},
		{
			email:         "",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			scenario: "Empty Email Input",
		},
		{
			email:         "invalid-email-format",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs("invalid-email-format").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			scenario: "Invalid Email Format",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			test.setupMocks(mock)

			result, err := s.GetByEmail(test.email)

			if test.expectedError != nil {
				if err == nil || err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error %v, but got %v", test.expectedError, err)
				}
				t.Logf("Scenario: %s succeeded: Correct error returned.", test.scenario)
			} else if err != nil {
				t.Errorf("Unexpected error occurred: %v", err)
			}

			if test.expectedUser != nil {
				if *result != *test.expectedUser {
					t.Errorf("Expected user %v, but got %v", *test.expectedUser, *result)
				}
				t.Logf("Scenario: %s succeeded: Correct user returned.", test.scenario)
			} else if result != nil {
				t.Errorf("Expected user to be nil, but got %v", result)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %s", err)
	}
}

/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1
*/
func TestGetByID(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		id        uint
		wantUser  *model.User
		wantErr   error
	}{
		{
			name: "Successful Retrieval of User by Valid ID",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(1).WillReturnRows(rows)
			},
			id:       1,
			wantUser: &model.User{ID: 1, Username: "testuser", Email: "test@example.com"},
			wantErr:  nil,
		},
		{
			name: "User Not Found in Database",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(2).WillReturnError(gorm.ErrRecordNotFound)
			},
			id:       2,
			wantUser: nil,
			wantErr:  gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(3).WillReturnError(errors.New("connection error"))
			},
			id:       3,
			wantUser: nil,
			wantErr:  errors.New("connection error"),
		},
		{
			name: "Retrieval with Invalid ID Input (Zero ID)",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(0).WillReturnError(errors.New("invalid query condition"))
			},
			id:       0,
			wantUser: nil,
			wantErr:  errors.New("invalid query condition"),
		},
		{
			name: "Retrieval with Invalid ID Input (Negative ID)",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE \"users\"\\.\"id\" = \\$1").
					WithArgs(-1).WillReturnError(errors.New("invalid query condition"))
			},
			id:       -1,
			wantUser: nil,
			wantErr:  errors.New("invalid query condition"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setupMock(mock)

			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			assert.NoError(t, err)

			userStore := &UserStore{db: gdb}

			gotUser, gotErr := userStore.GetByID(tt.id)

			if tt.wantUser != nil && gotUser != nil {
				assert.Equal(t, *tt.wantUser, *gotUser)
			} else {
				assert.Equal(t, tt.wantUser, gotUser)
			}

			assert.Equal(t, tt.wantErr, gotErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24
*/
func TestGetByUsername(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(mock sqlmock.Sqlmock)
		username      string
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "Retrieve User Successfully",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `username` = \\?").
					WithArgs("existing_user").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "existing_user"))
			},
			username:      "existing_user",
			expectedUser:  &model.User{ID: 1, Username: "existing_user"},
			expectedError: nil,
		},
		{
			name: "User Not Found",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `username` = \\?").
					WithArgs("nonexistent_user").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			username:      "nonexistent_user",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `username` = \\?").
					WithArgs("any_user").
					WillReturnError(errors.New("connection error"))
			},
			username:      "any_user",
			expectedUser:  nil,
			expectedError: errors.New("connection error"),
		},
		{
			name: "Username Case Sensitivity",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `username` = \\?").
					WithArgs("DifferentCase").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			username:      "DifferentCase",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Empty Username Input",
			setup: func(mock sqlmock.Sqlmock) {},
			username:      "",
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			assert.NoError(t, err)

			tt.setup(mock)

			store := &UserStore{db: gormDB}

			user, err := store.GetByUsername(tt.username)

			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedError, err)

			assert.NoError(t, mock.ExpectationsWereMet())

			if err != nil {
				t.Logf("Test '%s' failed with error: %v", tt.name, err)
			} else {
				t.Logf("Test '%s' passed successfully.", tt.name)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c
*/
func TestIsFollowing(t *testing.T) {
	type testCase struct {
		name      string
		userA     *model.User
		userB     *model.User
		mock      func(sqlmock.Sqlmock)
		expected  bool
		expectErr bool
	}

	tests := []testCase{
		{
			name:  "A follows B",
			userA: &model.User{ID: 1},
			userB: &model.User{ID: 2},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM follows`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:  true,
			expectErr: false,
		},
		{
			name:  "A does not follow B",
			userA: &model.User{ID: 1},
			userB: &model.User{ID: 2},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM follows`).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:  false,
			expectErr: false,
		},
		{
			name:  "A is nil, B is valid",
			userA: nil,
			userB: &model.User{ID: 2},
			mock: func(mock sqlmock.Sqlmock) {},
			expected:  false,
			expectErr: false,
		},
		{
			name:  "A is valid, B is nil",
			userA: &model.User{ID: 1},
			userB: nil,
			mock: func(mock sqlmock.Sqlmock) {},
			expected:  false,
			expectErr: false,
		},
		{
			name:  "Both A and B nil",
			userA: nil,
			userB: nil,
			mock: func(mock sqlmock.Sqlmock) {},
			expected:  false,
			expectErr: false,
		},
		{
			name:  "Database error occurs",
			userA: &model.User{ID: 1},
			userB: &model.User{ID: 2},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM follows`).
					WithArgs(1, 2).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			expected:  false,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			assert.NoError(t, err)

			userStore := &UserStore{db: gormDB}
			tt.mock(mock)

			got, err := userStore.IsFollowing(tt.userA, tt.userB)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}

			assert.Equal(t, tt.expected, got)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9
*/
func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name        string
		dbSetup     func() *gorm.DB
		expectedDB  *gorm.DB
		expectError bool
		logMessage  string
	}{
		{
			name: "Valid Database Connection",
			dbSetup: func() *gorm.DB {
				sqlDB, _, err := sqlmock.New()
				assert.NoError(t, err)

				gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
				assert.NoError(t, err)
				return gormDB
			},
			expectedDB:  nil,
			expectError: false,
			logMessage:  "Creating UserStore with a valid database connection should initialize correctly.",
		},
		{
			name: "Nil Database Connection",
			dbSetup: func() *gorm.DB {
				return nil
			},
			expectedDB:  nil,
			expectError: false,
			logMessage:  "Creating UserStore with a nil database connection should handle gracefully without errors.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.dbSetup()
			userStore := NewUserStore(db)

			if userStore.db == nil && tt.expectedDB == nil {
				t.Log(tt.logMessage)
			} else {
				assert.NotNil(t, userStore.db, "Expected non-nil DB connection")
				t.Log("UserStore db initialized as expected.")
			}

			expectedFields := 1
			actualFields := 1

			assert.Equal(t, expectedFields, actualFields, "UserStore should not have unexpected fields initialized")
			t.Log("UserStore initialization integrity maintained.")
		})
	}
}

/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55
*/
func TestUnfollow(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		prepare     func()
		userA       *model.User
		userB       *model.User
		expectedErr error
	}{
		{
			name: "Successfully Unfollowing a User",
			prepare: func() {
				mock.ExpectExec(`DELETE FROM "follows"`).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			userA:       &model.User{ID: 1},
			userB:       &model.User{ID: 2},
			expectedErr: nil,
		},
		{
			name: "Unfollowing a User Not Followed",
			prepare: func() {
				mock.ExpectExec(`DELETE FROM "follows"`).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			userA:       &model.User{ID: 1},
			userB:       &model.User{ID: 2},
			expectedErr: nil,
		},
		{
			name: "Database Error on Unfollow",
			prepare: func() {
				mock.ExpectExec(`DELETE FROM "follows"`).
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
			},
			userA:       &model.User{ID: 1},
			userB:       &model.User{ID: 2},
			expectedErr: errors.New("database error"),
		},
		{
			name: "Unfollowing with Nil User Parameters",
			prepare: func() {},
			userA:    nil,
			userB:    nil,
			expectedErr:  errors.New("invalid input"),
		},
		{
			name: "Unfollowing Same User",
			prepare: func() {
				mock.ExpectExec(`DELETE FROM "follows"`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			userA:       &model.User{ID: 1},
			userB:       &model.User{ID: 1},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			err = userStore.Unfollow(tt.userA, tt.userB)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435
*/
func TestUpdate(t *testing.T) {
	tests := []struct {
		name      string
		user      *model.User
		mockSetup func(sqlmock.Sqlmock)
		expectErr bool
		errType   error
	}{
		{
			name: "Successfully update a user in the database",
			user: &model.User{ID: 1, Username: "updated_user"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs("updated_user", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "Fail to update a user due to database error",
			user: &model.User{ID: 1, Username: "user_error"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs("user_error", 1).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectErr: true,
			errType:   errors.New("database error"),
		},
		{
			name: "Attempt to update a non-existent user",
			user: &model.User{ID: 999, Username: "non_existent"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs("non_existent", 999).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectErr: true,
			errType:   gorm.ErrRecordNotFound,
		},
		{
			name: "Attempt to update with invalid user data",
			user: &model.User{ID: 1, Username: ""},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			expectErr: true,
			errType:   errors.New("validation error"),
		},
		{
			name: "Concurrent update scenarios",
			user: &model.User{ID: 1, Username: "concurrent_user"},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WithArgs("concurrent_user", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
			assert.NoError(t, err)

			store := &UserStore{db: gdb}
			if tt.name == "Concurrent update scenarios" {
				var g errgroup.Group
				for i := 0; i < 5; i++ {
					g.Go(func() error {
						tt.mockSetup(mock)
						return store.Update(tt.user)
					})
				}

				if err := g.Wait(); err != nil && !tt.expectErr {
					t.Errorf("Unexpected error in concurrent test: %v", err)
				}
			} else {
				tt.mockSetup(mock)
				err := store.Update(tt.user)

				if !tt.expectErr && err != nil {
					t.Errorf("Expected no error, but got %v", err)
				}

				if tt.expectErr && !errors.Is(err, tt.errType) {
					t.Errorf("Expected error type %v, but got %v", tt.errType, err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Not all expectations were met: %v", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7
*/
func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	userStore := &UserStore{db: gormDB}

	type testData struct {
		description  string
		setup        func()
		input        *model.User
		expectedIDs  []uint
		expectingErr bool
	}

	tests := []testData{
		{
			description: "Successfully Retrieve Following User IDs",
			setup: func() {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)

				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE (.+)$").
					WithArgs(1).
					WillReturnRows(rows)
			},
			input:        &model.User{ID: 1},
			expectedIDs:  []uint{2, 3, 4},
			expectingErr: false,
		},
		{
			description: "No Following Users",
			setup: func() {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE (.+)$").
					WithArgs(2).
					WillReturnRows(rows)
			},
			input:        &model.User{ID: 2},
			expectedIDs:  []uint{},
			expectingErr: false,
		},
		{
			description: "Database Connection Error",
			setup: func() {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE (.+)$").
					WithArgs(3).
					WillReturnError(errors.New("database connection error"))
			},
			input:        &model.User{ID: 3},
			expectedIDs:  []uint{},
			expectingErr: true,
		},
		{
			description: "Invalid User Input (Nil User Object)",
			setup: func() {},
			input:        nil,
			expectedIDs:  []uint{},
			expectingErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.setup()

			ids, err := userStore.GetFollowingUserIDs(tc.input)
			assert.Equal(t, tc.expectingErr, err != nil, "Error mismatch")
			assert.ElementsMatch(t, tc.expectedIDs, ids, "IDs mismatch")

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
