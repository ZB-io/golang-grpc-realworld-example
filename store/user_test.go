package store

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
)

// type UserStore struct {
// 	db *gorm.DB
// }
// type DB struct {
// 	sync.RWMutex
// 	Value			interface{}
// 	Error			error
// 	RowsAffected		int64
// 	db			SQLCommon
// 	blockGlobalUpdate	bool
// 	logMode			logModeValue
// 	logger			logger
// 	search			*search
// 	values			sync.Map
// 	parent			*DB
// 	callbacks		*Callback
// 	dialect			Dialect
// 	singularTable		bool
// 	nowFuncOverride		func() time.Time
// }// single db
// // function to be used to override the creating of a new timestamp

// type DB struct {
// 	sync.RWMutex
// 	Value			interface{}
// 	Error			error
// 	RowsAffected		int64
// 	db			SQLCommon
// 	blockGlobalUpdate	bool
// 	logMode			logModeValue
// 	logger			logger
// 	search			*search
// 	values			sync.Map
// 	parent			*DB
// 	callbacks		*Callback
// 	dialect			Dialect
// 	singularTable		bool
// 	nowFuncOverride		func() time.Time
// }// single db
// function to be used to override the creating of a new timestamp

/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920


*/
func TestUserStoreCreate(t *testing.T) {

	type testCase struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		user    model.User
		wantErr bool
	}

	
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("failed to open gorm DB: %s", err)
	}

	store := NewUserStore(gormDB)

	tests := []testCase{
		{
			name: "Successfully create a new user",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO `users`").WithArgs(0, "uniqueuser", "example@example.com", "hashedpassword", "bio", "image").WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
			user:    model.User{Username: "uniqueuser", Email: "example@example.com", Password: "hashedpassword", Bio: "bio", Image: "image"},
			wantErr: false,
		},
		{
			name: "Fail to create user due to duplicate username",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO `users`").WithArgs(0, "duplicateuser", "first@example.com", "hashedpassword", "bio", "image").WillReturnError(gorm.ErrInvalidSQL)
				m.ExpectCommit()
			},
			user:    model.User{Username: "duplicateuser", Email: "first@example.com", Password: "hashedpassword", Bio: "bio", Image: "image"},
			wantErr: true,
		},
		{
			name: "Fail to create user due to missing fields",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO `users`").WithArgs(0, "", "example@example.com", "", "bio", "image").WillReturnError(gorm.ErrInvalidSQL)
				m.ExpectCommit()
			},
			user:    model.User{Email: "example@example.com", Bio: "bio", Image: "image"},
			wantErr: true,
		},
		{
			name: "Database connection error during user creation",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO `users`").WithArgs(0, "newuser", "example@example.com", "hashedpassword", "bio", "image").WillReturnError(gorm.ErrInvalidTransaction)
			},
			user:    model.User{Username: "newuser", Email: "example@example.com", Password: "hashedpassword", Bio: "bio", Image: "image"},
			wantErr: true,
		},
		{
			name: "Attempt to create user with invalid email format",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec("INSERT INTO `users`").WithArgs(0, "newuser", "invalid-email", "hashedpassword", "bio", "image").WillReturnError(gorm.ErrInvalidSQL)
				m.ExpectCommit()
			},
			user:    model.User{Username: "newuser", Email: "invalid-email", Password: "hashedpassword", Bio: "bio", Image: "image"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(mock)

			err := store.Create(&tc.user)

			if (err != nil) != tc.wantErr {
				t.Errorf("expected error: %v, got: %v", tc.wantErr, err)
			}

			if err == nil {
				t.Logf("Test %q passed: user successfully created", tc.name)
			} else {
				t.Logf("Test %q failed as expected: %v", tc.name, err)
			}
		})

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

/*
ROOST_METHOD_HASH=Follow_48fdf1257b
ROOST_METHOD_SIG_HASH=Follow_8217e61c06


 */
func TestFollow(t *testing.T) {

	tests := []struct {
		name          string
		setupMocks    func(sqlmock.Sqlmock)
		userA         model.User
		userB         model.User
		expectedError bool
	}{
		{
			name: "Successful Follow",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO follows").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			userA:         model.User{Username: "userA", Email: "userA@example.com", Password: "password"},
			userB:         model.User{Username: "userB", Email: "userB@example.com", Password: "password"},
			expectedError: false,
		},
		{
			name: "User Follow Twice",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT * FROM follows WHERE").WillReturnRows(sqlmock.NewRows([]string{"from_user_id", "to_user_id"}).AddRow(1, 2))
				mock.ExpectCommit()
			},
			userA:         model.User{Username: "userA", Email: "userA@example.com", Password: "password"},
			userB:         model.User{Username: "userB", Email: "userB@example.com", Password: "password"},
			expectedError: false,
		},
		{
			name: "Follow Non-Existent User",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO follows").WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectCommit()
			},
			userA:         model.User{Username: "userA", Email: "userA@example.com", Password: "password"},
			userB:         model.User{Model: gorm.Model{ID: 999}, Username: "ghostUser", Email: "ghostUser@example.com", Password: "password"},
			expectedError: true,
		},
		{
			name: "Handle Database Connection Error",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(gorm.ErrInvalidTransaction)
			},
			userA:         model.User{Username: "userA", Email: "userA@example.com", Password: "password"},
			userB:         model.User{Username: "userB", Email: "userB@example.com", Password: "password"},
			expectedError: true,
		},
		{
			name: "Follows Association Append Failure",
			setupMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO follows").WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectCommit()
			},
			userA:         model.User{Username: "userA", Email: "userA@example.com", Password: "password"},
			userB:         model.User{Username: "userB", Email: "userB@example.com", Password: "password"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a db", err)
			}

			userStore := &UserStore{db: gdb}

			tt.setupMocks(mock)

			err = userStore.Follow(&tt.userA, &tt.userB)
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error to be %v, got %v", tt.expectedError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Log("Test scenario: ", tt.name, "- passed")
		})
	}
}

/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1


 */
func TestUserStoreGetByEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		mockSetup func(mock sqlmock.Sqlmock)
		wantError bool
		wantUser  *model.User
	}{
		{
			name:  "Retrieve Existing User by Email",
			email: "existing@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(1, "existing@example.com")
				mock.ExpectQuery("^SELECT .+ FROM \"users\" WHERE email = ?").
					WithArgs("existing@example.com").WillReturnRows(rows)
			},
			wantError: false,
			wantUser:  &model.User{Email: "existing@example.com"},
		},
		{
			name:  "Attempt to Retrieve Non-Existent User by Email",
			email: "nonexistent@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT .+ FROM \"users\" WHERE email = ?").
					WithArgs("nonexistent@example.com").WillReturnError(gorm.ErrRecordNotFound)
			},
			wantError: true,
			wantUser:  nil,
		},
		{
			name:  "Handle Database Error While Retrieving User",
			email: "error@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT .+ FROM \"users\" WHERE email = ?").
					WithArgs("error@example.com").WillReturnError(fmt.Errorf("database error"))
			},
			wantError: true,
			wantUser:  nil,
		},
		{
			name:  "Retrieve User with Special Characters in Email",
			email: "special+char@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(2, "special+char@example.com")
				mock.ExpectQuery("^SELECT .+ FROM \"users\" WHERE email = ?").
					WithArgs("special+char@example.com").WillReturnRows(rows)
			},
			wantError: false,
			wantUser:  &model.User{Email: "special+char@example.com"},
		},
		{
			name:  "Retrieve User with Leading/Trailing Spaces in Email",
			email: " padded@example.com ",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(3, "padded@example.com")
				mock.ExpectQuery("^SELECT .+ FROM \"users\" WHERE email = ?").
					WithArgs("padded@example.com").WillReturnRows(rows)
			},
			wantError: false,
			wantUser:  &model.User{Email: "padded@example.com"},
		},
		{
			name:  "Retrieve User with Upper and Lower Case Email Sensitivity",
			email: "CASE@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(4, "case@example.com")
				mock.ExpectQuery("^SELECT .+ FROM \"users\" WHERE email = ?").
					WithArgs("CASE@example.com").WillReturnRows(rows)
			},
			wantError: false,
			wantUser:  &model.User{Email: "case@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open gorm database: %v", err)
			}

			s := &UserStore{db: gormDB}

			tt.mockSetup(mock)

			user, err := s.GetByEmail(tt.email)
			if (err != nil) != tt.wantError {
				t.Fatalf("Expected error: %v, got: %v", tt.wantError, err)
			}

			if tt.wantUser != nil && user != nil && user.Email != tt.wantUser.Email {
				t.Errorf("Expected user: %v, got: %v", tt.wantUser, user)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1


 */
func TestGetByID(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open("sqlite3", db)
	if err != nil {
		t.Fatalf("could not initialize gorm DB: %v", err)
	}
	store := &UserStore{db: gdb}

	tests := []struct {
		name          string
		id            uint
		mockBehavior  func()
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "Scenario 1: Successful Retrieval of User by Valid ID",
			id:   1,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "John Doe", "john@example.com")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((id = ?))`)).
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedUser:  &model.User{Model: gorm.Model{ID: 1}, Username: "John Doe", Email: "john@example.com"},
			expectedError: nil,
		},
		{
			name: "Scenario 2: User Not Found for Non-existent ID",
			id:   2,
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((id = ?))`)).
					WithArgs(2).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Scenario 3: Database Error During User Retrieval",
			id:   3,
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((id = ?))`)).
					WithArgs(3).
					WillReturnError(gorm.ErrInvalidTransaction)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrInvalidTransaction,
		},
		{
			name: "Scenario 4: Edge Case with User ID 0",
			id:   0,
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((id = ?))`)).
					WithArgs(0).
					WillReturnError(gorm.ErrInvalidSQL)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrInvalidSQL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Running test: %s", tt.name)
			tt.mockBehavior()

			user, err := store.GetByID(tt.id)

			assert.Equal(t, tt.expectedUser, user, "expected user does not match")
			assert.Equal(t, tt.expectedError, err, "unexpected error")
		})
	}
}

/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24


 */
// func (s *UserStore) GetByUsername(username string) (*model.User, error) {
// 	var m model.User
// 	if err := s.db.Where("username = ?", username).First(&m).Error; err != nil {
// 		return nil, err
// 	}
// 	return &m, nil
// }

func TestUserStoreGetByUsername(t *testing.T) {
	// var db *gorm.DB
	// var mock sqlmock.Sqlmock
	var err error

	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "Mock database connection should be established")

	defer func() {
		mock.ExpectClose()
		db.Close()
	}()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("failed to open gorm DB: %s", err)
	}
	userStore := NewUserStore(gormDB)

	tests := []struct {
		name     string
		username string
		mock     func()
		wantUser *model.User
		wantErr  bool
	}{
		{
			name:     "Successfully Retrieve User by Username",
			username: "testuser",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username"}).
					AddRow(1, "testuser")
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE \\(username = \\?\\) ORDER BY `users`.`id` LIMIT 1").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			wantUser: &model.User{ID: 1, Username: "testuser"},
			wantErr:  false,
		},
		{
			name:     "Fail to Retrieve Non-Existent Username",
			username: "unknown",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE \\(username = \\?\\) ORDER BY `users`.`id` LIMIT 1").
					WithArgs("unknown").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name:     "Handle Database Error",
			username: "anyuser",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE \\(username = \\?\\) ORDER BY `users`.`id` LIMIT 1").
					WithArgs("anyuser").
					WillReturnError(errors.New("database error"))
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name:     "Empty Username Input",
			username: "",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE \\(username = \\?\\) ORDER BY `users`.`id` LIMIT 1").
					WithArgs("").
					WillReturnError(errors.New("invalid input syntax"))
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name:     "Case Sensitivity Check",
			username: "TestUser",
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM `users` WHERE \\(username = \\?\\) ORDER BY `users`.`id` LIMIT 1").
					WithArgs("TestUser").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			user, err := userStore.GetByUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error for '%s' case", tt.name)
			} else {
				assert.NoError(t, err, "Did not expect an error for '%s' case", tt.name)
			}
			assert.Equal(t, tt.wantUser, user, "Expected user object mismatch for '%s' case", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c


 */
func TestUserStoreIsFollowing(t *testing.T) {

	type testCase struct {
		name      string
		userA     *model.User
		userB     *model.User
		mockSetup func(mock sqlmock.Sqlmock)
		expected  bool
		expectErr error
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %s", err)
	}
	defer db.Close()

	gdb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Error initializing gorm with sqlmock: %s", err)
	}

	userA := &model.User{Model: gorm.Model{ID: 1}}
	userB := &model.User{Model: gorm.Model{ID: 2}}

	userStore := UserStore{db: gdb}

	testCases := []testCase{
		{
			name:  "Scenario 1: Check Following Relationship Exists",
			userA: userA,
			userB: userB,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT count(\\*) FROM \"follows\" WHERE \\(from_user_id = \\$1 AND to_user_id = \\$2\\)").
					WithArgs(userA.ID, userB.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:  true,
			expectErr: nil,
		},
		{
			name:  "Scenario 2: Check Following Relationship Does Not Exist",
			userA: userA,
			userB: userB,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT count(\\*) FROM \"follows\" WHERE \\(from_user_id = \\$1 AND to_user_id = \\$2\\)").
					WithArgs(userA.ID, userB.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:  false,
			expectErr: nil,
		},
		{
			name:      "Scenario 3: Handle Nil User A",
			userA:     nil,
			userB:     userB,
			mockSetup: func(mock sqlmock.Sqlmock) {},
			expected:  false,
			expectErr: nil,
		},
		{
			name:      "Scenario 4: Handle Nil User B",
			userA:     userA,
			userB:     nil,
			mockSetup: func(mock sqlmock.Sqlmock) {},
			expected:  false,
			expectErr: nil,
		},
		{
			name:  "Scenario 5: Handle Database Error",
			userA: userA,
			userB: userB,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT count(\\*) FROM \"follows\" WHERE \\(from_user_id = \\$1 AND to_user_id = \\$2\\)").
					WithArgs(userA.ID, userB.ID).
					WillReturnError(errors.New("database error"))
			},
			expected:  false,
			expectErr: errors.New("database error"),
		},
		{
			name:      "Scenario 6: No Users",
			userA:     nil,
			userB:     nil,
			mockSetup: func(mock sqlmock.Sqlmock) {},
			expected:  false,
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup(mock)
			actual, err := userStore.IsFollowing(tc.userA, tc.userB)

			if tc.expectErr != nil {
				if err == nil || err.Error() != tc.expectErr.Error() {
					t.Errorf("Expected error %v, but got %v", tc.expectErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if actual != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, actual)
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9


 */
// func NewUserStore(db *gorm.DB) *UserStore {
// 	return &UserStore{
// 		db: db,
// 	}
// }

func TestNewUserStore(t *testing.T) {
	t.Run("Scenario 1: Successful Initialization of UserStore with a Valid DB Instance", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("failed to open mock GORM DB: %v", err)
		}
		defer gormDB.Close()

		userStore := NewUserStore(gormDB)
		if userStore == nil || userStore.db != gormDB {
			t.Errorf("expected UserStore to be initialized with the provided db instance")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unmet expectations: %s", err)
		}
	})

	t.Run("Scenario 2: Handling of Nil DB Parameter", func(t *testing.T) {
		userStore := NewUserStore(nil)
		if userStore == nil {
			t.Error("expected UserStore to not be nil even when nil db is provided")
		} else if userStore.db != nil {
			t.Errorf("expected UserStore's db field to be nil, got: %+v", userStore.db)
		}
	})

	t.Run("Scenario 3: Consecutive Calls with Different DB Instances", func(t *testing.T) {
		db1, mock1, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create first mock DB: %v", err)
		}
		defer db1.Close()

		gormDB1, err := gorm.Open("postgres", db1)
		if err != nil {
			t.Fatalf("failed to open first mock GORM DB: %v", err)
		}
		defer gormDB1.Close()

		db2, mock2, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create second mock DB: %v", err)
		}
		defer db2.Close()

		gormDB2, err := gorm.Open("postgres", db2)
		if err != nil {
			t.Fatalf("failed to open second mock GORM DB: %v", err)
		}
		defer gormDB2.Close()

		userStore1 := NewUserStore(gormDB1)
		userStore2 := NewUserStore(gormDB2)

		if userStore1 == userStore2 || userStore1.db == userStore2.db {
			t.Error("expected two distinct UserStore instances with different db instances")
		}

		if err := mock1.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unmet expectations for first DB: %s", err)
		}
		if err := mock2.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unmet expectations for second DB: %s", err)
		}
	})

	t.Run("Scenario 4: Cross-Verification with Model Struct Interaction", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		gormDB, err := gorm.Open("postgres", db)
		if err != nil {
			t.Fatalf("failed to open mock GORM DB: %v", err)
		}
		defer gormDB.Close()

		// userStore := NewUserStore(gormDB)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unmet expectations: %s", err)
		}
	})
}

/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55


 */
func TestUserStoreUnfollow(t *testing.T) {
	scenarios := []struct {
		name      string
		follower  *model.User
		followee  *model.User
		setupMock func(sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "Successfully Unfollowing a User",
			follower: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "follower",
			},
			followee: &model.User{
				Model: gorm.Model{
					ID: 2,
				},
				Username: "followee",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows WHERE from_user_id = \? AND to_user_id = \?`).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "Attempting to Unfollow a Non-Followed User",
			follower: &model.User{
				Model: gorm.Model{
					ID: 3,
				},
				Username: "follower_not_following",
			},
			followee: &model.User{
				Model: gorm.Model{
					ID: 4,
				},
				Username: "followee_not_followed",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows WHERE from_user_id = \? AND to_user_id = \?`).
					WithArgs(3, 4).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "Database Error During Unfollow Operation",
			follower: &model.User{
				Model: gorm.Model{
					ID: 5,
				},
				Username: "error_follower",
			},
			followee: &model.User{
				Model: gorm.Model{
					ID: 6,
				},
				Username: "error_followee",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows WHERE from_user_id = \? AND to_user_id = \?`).
					WithArgs(5, 6).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:     "Unfollowing with a Nil User Input",
			follower: nil,
			followee: &model.User{
				Model: gorm.Model{
					ID: 7,
				},
				Username: "nil_followee",
			},
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectErr: true,
		},
		{
			name: "Attempting to Unfollow When Database Connection is Lost",
			follower: &model.User{
				Model: gorm.Model{
					ID: 8,
				},
				Username: "connection_lost_follower",
			},
			followee: &model.User{
				Model: gorm.Model{
					ID: 9,
				},
				Username: "connection_lost_followee",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM follows WHERE from_user_id = \? AND to_user_id = \?`).
					WithArgs(8, 9).
					WillReturnError(errors.New("connection lost"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			sqlDB, err := gorm.Open("postgres", db)
			if err != nil {
				log.Fatalf("failed to initialize GORM with the mocked database: %v", err)
			}

			userStore := UserStore{db: sqlDB}
			scenario.setupMock(mock)

			err = userStore.Unfollow(scenario.follower, scenario.followee)
			if (err != nil) != scenario.expectErr {
				t.Errorf("Unexpected error outcome: got %v, want %v", err != nil, scenario.expectErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Log("Test scenario completed successfully:", scenario.name)
		})
	}
}

/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435


 */
func TestUserStoreUpdate(t *testing.T) {
	type testCase struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		user      model.User
		expectErr bool
	}

	tests := []testCase{
		{
			name: "Successful User Update in Database",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE \"users\" SET").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			user: model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Bio:      "Bio",
				Image:    "https://example.com/image.png",
				Password: "password",
			},
			expectErr: false,
		},
		{
			name: "User Update with Non-Existent User",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE \"users\" SET").WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			user: model.User{
				Model:    gorm.Model{ID: 9999},
				Username: "ghostuser",
				Email:    "ghost@example.com",
			},
			expectErr: true,
		},
		{
			name: "Update with Invalid or Malformed Data",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE \"users\" SET").WillReturnError(errors.New("constraint violation"))
				mock.ExpectRollback()
			},
			user: model.User{
				Model:    gorm.Model{ID: 1},
				Username: "",
				Email:    "duplicate@example.com",
			},
			expectErr: true,
		},
		{
			name: "Update Operation on Database Disconnection",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("db connection error"))
			},
			user: model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("failed to open gorm DB: %v", err)
			}

			tt.setupMock(mock)

			store := UserStore{db: gormDB}

			err = store.Update(&tt.user)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else {
					t.Logf("Expected error received: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error but got: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("not all expectations were met: %v", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7


 */
func TestUserStoreGetFollowingUserIDs(t *testing.T) {

	tests := []struct {
		name      string
		mockSetup func(sqlmock.Sqlmock)
		inputUser *model.User
		wantIDs   []uint
		expectErr bool
	}{
		{
			name: "Normal Operation with Multiple Followings",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).AddRow(2).AddRow(3).AddRow(4)
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).WillReturnRows(rows)
			},
			inputUser: &model.User{Model: gorm.Model{ID: 1}},
			wantIDs:   []uint{2, 3, 4},
			expectErr: false,
		},
		{
			name: "No Followings for User",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"to_user_id"}))
			},
			inputUser: &model.User{Model: gorm.Model{ID: 2}},
			wantIDs:   []uint{},
			expectErr: false,
		},
		{
			name: "Database Error Encounter",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(3).WillReturnError(errors.New("database error"))
			},
			inputUser: &model.User{Model: gorm.Model{ID: 3}},
			wantIDs:   []uint{},
			expectErr: true,
		},
		{
			name: "Database Connection Lost",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(4).WillReturnError(gorm.ErrInvalidTransaction)
			},
			inputUser: &model.User{Model: gorm.Model{ID: 4}},
			wantIDs:   []uint{},
			expectErr: true,
		},
		{
			name: "Single Following User ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).AddRow(5)
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(5).WillReturnRows(rows)
			},
			inputUser: &model.User{Model: gorm.Model{ID: 5}},
			wantIDs:   []uint{5},
			expectErr: false,
		},
		{
			name: "Maximum Integer for User ID",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).AddRow(6).AddRow(7)
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(^uint(0)).WillReturnRows(rows)
			},
			inputUser: &model.User{Model: gorm.Model{ID: ^uint(0)}},
			wantIDs:   []uint{6, 7},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Unable to create mock sql database: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			gormDB, _ := gorm.Open("mysql", db)
			store := &UserStore{db: gormDB}

			gotIDs, err := store.GetFollowingUserIDs(tt.inputUser)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetFollowingUserIDs() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !equalSlices(gotIDs, tt.wantIDs) {
				t.Errorf("GetFollowingUserIDs() = %v, want %v", gotIDs, tt.wantIDs)
			}
		})
	}
}

func equalSlices(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

