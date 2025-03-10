package store

import (
	debug "runtime/debug"
	testing "testing"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	sql "database/sql"
	regexp "regexp"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	fmt "fmt"
	math "math"
	time "time"
	assert "github.com/stretchr/testify/assert"
	errors "errors"
)








/*
ROOST_METHOD_HASH=NewUserStore_fb599438e5
ROOST_METHOD_SIG_HASH=NewUserStore_c0075221af

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore // NewUserStore returns a new UserStore


*/
func TestNewUserStore(t *testing.T) {

	type testCase struct {
		name      string
		db        *gorm.DB
		wantNil   bool
		setupMock func(sqlmock.Sqlmock)
	}

	tests := []testCase{
		{
			name: "Successfully Create New UserStore Instance",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.0"))
			},
			wantNil: false,
		},
		{
			name:      "Create UserStore with Nil Database Connection",
			db:        nil,
			wantNil:   true,
			setupMock: func(mock sqlmock.Sqlmock) {},
		},
		{
			name: "Verify UserStore Instance Independence",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.0"))
			},
			wantNil: false,
		},
		{
			name: "Create UserStore with Configured Database Connection",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.0"))
			},
			wantNil: false,
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

			var mockDB *gorm.DB
			if tt.db == nil && !tt.wantNil {
				sqlDB, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Failed to create mock: %v", err)
				}
				defer sqlDB.Close()

				tt.setupMock(mock)
				mockDB, err = gorm.Open("mysql", sqlDB)
				if err != nil {
					t.Fatalf("Failed to open gorm connection: %v", err)
				}
			} else {
				mockDB = tt.db
			}

			got := NewUserStore(mockDB)

			if tt.wantNil {
				if got.db != nil {
					t.Errorf("NewUserStore() = %v, want nil DB", got)
				}
				t.Log("Successfully verified nil DB handling")
				return
			}

			if got == nil {
				t.Error("NewUserStore() returned nil, want non-nil UserStore")
				return
			}

			if got.db == nil {
				t.Error("NewUserStore() returned UserStore with nil DB, want non-nil DB")
				return
			}

			if tt.name == "Verify UserStore Instance Independence" {
				store2 := NewUserStore(mockDB)
				if got == store2 {
					t.Error("NewUserStore() returned same instance for different calls")
				}
			}

			t.Logf("Successfully completed test case: %s", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_Create_9495ddb29d
ROOST_METHOD_SIG_HASH=UserStore_Create_18451817fe

FUNCTION_DEF=func (s *UserStore) Create(m *model.User) error // Create create a user


*/
func TestUserStoreCreate(t *testing.T) {

	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
		setup   func(mock gosqlmock.Sqlmock)
	}{
		{
			name: "Successful User Creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			wantErr: false,
			setup: func(mock gosqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WithArgs(gosqlmock.AnyArg(), gosqlmock.AnyArg(), gosqlmock.AnyArg(), "testuser", "test@example.com", "password123", "Test bio", "https://example.com/image.jpg").
					WillReturnResult(gosqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "Duplicate Username",
			user: &model.User{
				Username: "existinguser",
				Email:    "new@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			wantErr: true,
			setup: func(mock gosqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnError(&sql.Error{Message: "Duplicate entry 'existinguser' for key 'username'"})
				mock.ExpectRollback()
			},
		},
		{
			name: "Duplicate Email",
			user: &model.User{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			wantErr: true,
			setup: func(mock gosqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnError(&sql.Error{Message: "Duplicate entry 'existing@example.com' for key 'email'"})
				mock.ExpectRollback()
			},
		},
		{
			name: "Missing Required Fields",
			user: &model.User{
				Username: "",
				Email:    "",
				Password: "",
			},
			wantErr: true,
			setup: func(mock gosqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnError(&sql.Error{Message: "Column 'username' cannot be null"})
				mock.ExpectRollback()
			},
		},
		{
			name: "Database Connection Error",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
			setup: func(mock gosqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnError(&sql.Error{Message: "Connection refused"})
				mock.ExpectRollback()
			},
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

			db, mock, err := gosqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer gormDB.Close()

			tt.setup(mock)

			store := &UserStore{
				db: gormDB,
			}

			err = store.Create(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test '%s' failed with error: %v", tt.name, err)
			} else {
				t.Logf("Test '%s' passed successfully", tt.name)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_Follow_fe0976e4eb
ROOST_METHOD_SIG_HASH=UserStore_Follow_0e703b23f8

FUNCTION_DEF=func (s *UserStore) Follow(a *model.User, b *model.User) error // Follow create follow relashionship to User B from user A


*/
func TestUserStoreFollow(t *testing.T) {

	type testCase struct {
		name        string
		userA       *model.User
		userB       *model.User
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
	}

	tests := []testCase{
		{
			name: "Successful Follow",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Follow Already Followed User",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Self Follow Attempt",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 1).
					WillReturnError(fmt.Errorf("self-follow not allowed"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name:  "Follow with Nil User A",
			userA: nil,
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
		},
		{
			name: "Follow with Nil User B",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB:       nil,
			setupMock:   func(mock sqlmock.Sqlmock) {},
			expectError: true,
		},
		{
			name: "Database Error",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnError(fmt.Errorf("database error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.setupMock(mock)

			store := &UserStore{db: gormDB}

			err = store.Follow(tc.userA, tc.userB)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_GetByEmail_fda09af5c4
ROOST_METHOD_SIG_HASH=UserStore_GetByEmail_9e84f3286b

FUNCTION_DEF=func (s *UserStore) GetByEmail(email string) (*model.User, error) // GetByEmail finds a user from email


*/
func TestUserStoreGetByEmail(t *testing.T) {

	type testCase struct {
		name          string
		email         string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}

	tests := []testCase{
		{
			name:  "Successfully retrieve user by valid email",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg")
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
		{
			name:  "Non-existent email",
			email: "nonexistent@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("nonexistent@example.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Empty email string",
			email: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Database connection error",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &UserStore{
				db: gormDB,
			}

			user, err := store.GetByEmail(tc.email)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}

			if tc.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v but got nil", tc.expectedError)
				} else if err != tc.expectedError {
					t.Errorf("Expected error %v but got %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got %v", err)
				}
			}

			if tc.expectedUser != nil {
				if user == nil {
					t.Error("Expected user but got nil")
				} else {
					if user.Email != tc.expectedUser.Email {
						t.Errorf("Expected email %s but got %s", tc.expectedUser.Email, user.Email)
					}

				}
			} else if user != nil {
				t.Error("Expected nil user but got a user")
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_GetByID_1f5f06165b
ROOST_METHOD_SIG_HASH=UserStore_GetByID_2a864916bb

FUNCTION_DEF=func (s *UserStore) GetByID(id uint) (*model.User, error) // GetByID finds a user from id


*/
func TestUserStoreGetById(t *testing.T) {

	type testCase struct {
		name          string
		id            uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully retrieve existing user",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, "2023-01-01", "2023-01-01", nil, "testuser", "test@example.com", "hashedpass", "bio", "image.jpg"))
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpass",
				Bio:      "bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
		{
			name: "User not found",
			id:   999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database connection error",
			id:   1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name: "Zero ID handling",
			id:   0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Maximum uint value handling",
			id:   math.MaxUint32,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `users` WHERE `users`.`id` = ?").
					WithArgs(math.MaxUint32).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &UserStore{db: gormDB}

			t.Logf("Executing test case: %s", tc.name)
			user, err := store.GetByID(tc.id)

			if err != tc.expectedError {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}

			if tc.expectedUser != nil {
				if user == nil {
					t.Error("Expected user to not be nil")
				} else if user.ID != tc.expectedUser.ID {
					t.Errorf("Expected user ID %d, got %d", tc.expectedUser.ID, user.ID)
				}

			} else if user != nil {
				t.Error("Expected user to be nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case completed: %s", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_GetByUsername_622b1b9e41
ROOST_METHOD_SIG_HASH=UserStore_GetByUsername_992f00baec

FUNCTION_DEF=func (s *UserStore) GetByUsername(username string) (*model.User, error) // GetByUsername finds a user from username


*/
func TestUserStoreGetByUsername(t *testing.T) {

	type testCase struct {
		name          string
		username      string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []testCase{
		{
			name:     "Success - Valid Username",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(1, fixedTime, fixedTime, nil, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg")
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(username = \?\) ORDER BY "users"\."id" ASC LIMIT 1`).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
		{
			name:     "Failure - User Not Found",
			username: "nonexistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(username = \?\) ORDER BY "users"\."id" ASC LIMIT 1`).
					WithArgs("nonexistent").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Failure - Empty Username",
			username: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(username = \?\) ORDER BY "users"\."id" ASC LIMIT 1`).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Failure - Database Error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE \(username = \?\) ORDER BY "users"\."id" ASC LIMIT 1`).
					WithArgs("testuser").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &UserStore{db: gormDB}

			t.Logf("Executing test case: %s", tc.name)
			user, err := store.GetByUsername(tc.username)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, user)
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectedUser.Username, user.Username)
				assert.Equal(t, tc.expectedUser.Email, user.Email)
				t.Logf("Successfully retrieved user: %v", user.Username)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_GetFollowingUserIDs_ee9c1008ff
ROOST_METHOD_SIG_HASH=UserStore_GetFollowingUserIDs_d7746035ec

FUNCTION_DEF=func (s *UserStore) GetFollowingUserIDs(m *model.User) ([ // GetFollowingUserIDs returns user ids current user follows
]uint, error) 

*/
func TestUserStoreGetFollowingUserIDs(t *testing.T) {

	type testCase struct {
		name          string
		user          *model.User
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedIDs   []uint
		expectedError error
	}

	tests := []testCase{
		{
			name: "Successfully retrieve following user IDs",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(uint(2)).
					AddRow(uint(3)).
					AddRow(uint(4))
				mock.ExpectQuery("SELECT (.+) FROM `follows` WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: nil,
		},
		{
			name: "Empty following list",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT (.+) FROM `follows` WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name: "Database error",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `follows` WHERE").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			expectedIDs:   []uint{},
			expectedError: errors.New("database error"),
		},
		{
			name: "Invalid user ID",
			user: &model.User{Model: gorm.Model{ID: 0}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `follows` WHERE").
					WithArgs(0).
					WillReturnError(sql.ErrNoRows)
			},
			expectedIDs:   []uint{},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			store := &UserStore{db: gormDB}

			ids, err := store.GetFollowingUserIDs(tc.user)

			if err != nil && tc.expectedError == nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if err == nil && tc.expectedError != nil {
				t.Errorf("Expected error %v, got nil", tc.expectedError)
			}
			if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}

			if len(ids) != len(tc.expectedIDs) {
				t.Errorf("Expected %d IDs, got %d", len(tc.expectedIDs), len(ids))
			}
			for i, id := range ids {
				if id != tc.expectedIDs[i] {
					t.Errorf("Expected ID %d at position %d, got %d", tc.expectedIDs[i], i, id)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_IsFollowing_309135b43a
ROOST_METHOD_SIG_HASH=UserStore_IsFollowing_4644b1529c

FUNCTION_DEF=func (s *UserStore) IsFollowing(a *model.User, b *model.User) (bool, error) // IsFollowing returns whether user A follows user B or not


*/
func TestUserStoreIsFollowing(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name     string
		userA    *model.User
		userB    *model.User
		mockFunc func()
		want     bool
		wantErr  bool
	}{
		{
			name:  "Valid Following Relationship",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockFunc: func() {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want:    true,
			wantErr: false,
		},
		{
			name:  "Non-Existent Following Relationship",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockFunc: func() {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want:    false,
			wantErr: false,
		},
		{
			name:     "Nil User A",
			userA:    nil,
			userB:    &model.User{Model: gorm.Model{ID: 2}},
			mockFunc: func() {},
			want:     false,
			wantErr:  false,
		},
		{
			name:     "Nil User B",
			userA:    &model.User{Model: gorm.Model{ID: 1}},
			userB:    nil,
			mockFunc: func() {},
			want:     false,
			wantErr:  false,
		},
		{
			name:  "Database Error",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockFunc: func() {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(sql.ErrConnDone)
			},
			want:    false,
			wantErr: true,
		},
		{
			name:  "Same User Check",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 1}},
			mockFunc: func() {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want:    false,
			wantErr: false,
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

			t.Logf("Testing scenario: %s", tt.name)

			tt.mockFunc()

			got, err := userStore.IsFollowing(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.IsFollowing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserStore.IsFollowing() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test completed successfully")
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_Unfollow_29d3ef7f50
ROOST_METHOD_SIG_HASH=UserStore_Unfollow_31d9214353

FUNCTION_DEF=func (s *UserStore) Unfollow(a *model.User, b *model.User) error // Unfollow delete follow relashionship to User B from user A


*/
func TestUserStoreUnfollow(t *testing.T) {

	type testCase struct {
		name    string
		userA   *model.User
		userB   *model.User
		dbSetup func(mock sqlmock.Sqlmock)
		wantErr bool
	}

	validUserA := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "userA",
		Email:    "userA@test.com",
		Password: "password",
	}

	validUserB := &model.User{
		Model:    gorm.Model{ID: 2},
		Username: "userB",
		Email:    "userB@test.com",
		Password: "password",
	}

	tests := []testCase{
		{
			name:  "Successful unfollow",
			userA: validUserA,
			userB: validUserB,
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(validUserA.ID, validUserB.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "Unfollow non-existent relationship",
			userA: validUserA,
			userB: validUserB,
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(validUserA.ID, validUserB.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "Invalid user A (nil)",
			userA: nil,
			userB: validUserB,
			dbSetup: func(mock sqlmock.Sqlmock) {

			},
			wantErr: true,
		},
		{
			name:  "Invalid user B (nil)",
			userA: validUserA,
			userB: nil,
			dbSetup: func(mock sqlmock.Sqlmock) {

			},
			wantErr: true,
		},
		{
			name:  "Database error",
			userA: validUserA,
			userB: validUserB,
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(validUserA.ID, validUserB.ID).
					WillReturnError(gorm.ErrInvalidTransaction)
				mock.ExpectRollback()
			},
			wantErr: true,
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

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer gormDB.Close()

			if tt.dbSetup != nil {
				tt.dbSetup(mock)
			}

			store := &UserStore{
				db: gormDB,
			}

			err = store.Unfollow(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Unfollow() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_Update_4fd6d3d1c1
ROOST_METHOD_SIG_HASH=UserStore_Update_ddd5c151cf

FUNCTION_DEF=func (s *UserStore) Update(m *model.User) error // Update update all of user fields


*/
func TestUserStoreUpdate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name    string
		user    *model.User
		setupFn func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successful Update",
			user: &model.User{
				Model: gorm.Model{
					ID:        1,
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
				Bio:      "Updated bio",
				Image:    "updated-image.jpg",
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WithArgs(
						sqlmock.AnyArg(),
						"testuser",
						"test@example.com",
						sqlmock.AnyArg(),
						"Updated bio",
						"updated-image.jpg",
						1,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Update Non-Existent User",
			user: &model.User{
				Model: gorm.Model{
					ID: 999,
				},
				Username: "nonexistent",
				Email:    "nonexistent@example.com",
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(fmt.Errorf("record not found"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Update With Invalid Data",
			user: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "",
				Email:    "invalid-email",
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(fmt.Errorf("validation error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Update With Null Optional Fields",
			user: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "testuser",
				Email:    "test@example.com",
				Bio:      "",
				Image:    "",
			},
			setupFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WithArgs(
						sqlmock.AnyArg(),
						"testuser",
						"test@example.com",
						sqlmock.AnyArg(),
						"",
						"",
						1,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
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

			tt.setupFn(mock)

			err := userStore.Update(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err != nil {
				t.Logf("Test '%s' completed with expected error: %v", tt.name, err)
			} else {
				t.Logf("Test '%s' completed successfully", tt.name)
			}
		})
	}
}

