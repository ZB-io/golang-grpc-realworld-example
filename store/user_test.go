package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"database/sql"
	"math"
	"time"
	"regexp"
)

type T struct {
	common
	isEnvSet bool
	context  *testContext
}
type ExpectedBegin struct {
	commonExpectation
	delay time.Duration
}
type ExpectedCommit struct {
	commonExpectation
}
type ExpectedExec struct {
	queryBasedExpectation
	result driver.Result
	delay  time.Duration
}
type ExpectedRollback struct {
	commonExpectation
}
type ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}
type Rows struct {
	converter driver.ValueConverter
	cols      []string
	def       []*Column
	rows      [][]driver.Value
	pos       int
	nextErr   map[int]error
	closeErr  error
}
type Time struct {
	wall uint64
	ext  int64

	loc *Location
}
/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9


 */
func TestNewUserStore(t *testing.T) {

	type testCase struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	gormDB, err := gorm.Open("sqlite3", mockDB)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	mockDB2, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create second mock DB: %v", err)
	}
	defer mockDB2.Close()

	gormDB2, err := gorm.Open("sqlite3", mockDB2)
	if err != nil {
		t.Fatalf("Failed to open second gorm connection: %v", err)
	}
	defer gormDB2.Close()

	configuredDB := gormDB.LogMode(true)

	tests := []testCase{
		{
			name:     "Scenario 1: Successfully Create New UserStore with Valid DB Connection",
			db:       gormDB,
			wantNil:  false,
			scenario: "Basic initialization with valid DB",
		},
		{
			name:     "Scenario 2: Create UserStore with Nil DB Parameter",
			db:       nil,
			wantNil:  false,
			scenario: "Nil DB initialization",
		},
		{
			name:     "Scenario 3: Verify DB Reference Integrity",
			db:       gormDB,
			wantNil:  false,
			scenario: "DB reference verification",
		},
		{
			name:     "Scenario 4: Multiple UserStore Instances Independence",
			db:       gormDB2,
			wantNil:  false,
			scenario: "Multiple instance verification",
		},
		{
			name:     "Scenario 5: UserStore Creation with Configured DB Instance",
			db:       configuredDB,
			wantNil:  false,
			scenario: "Configured DB initialization",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Starting test:", tc.scenario)

			userStore := NewUserStore(tc.db)

			if tc.wantNil {
				assert.Nil(t, userStore, "Expected nil UserStore")
			} else {
				assert.NotNil(t, userStore, "Expected non-nil UserStore")
				assert.Equal(t, tc.db, userStore.db, "DB reference mismatch")
			}

			switch tc.scenario {
			case "Multiple instance verification":

				userStore2 := NewUserStore(gormDB)
				assert.NotEqual(t, userStore.db, userStore2.db, "UserStore instances should have independent DB references")

			case "Configured DB initialization":
				logMode := userStore.db.LogMode(true)
				assert.NotNil(t, logMode, "DB configuration not preserved")
			}

			t.Log("Test completed successfully:", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920


 */
func TestUserStoreCreate(t *testing.T) {

	tests := []struct {
		name    string
		user    *model.User
		dbError error
		wantErr bool
	}{
		{
			name: "Successful user creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Duplicate username",
			user: &model.User{
				Username: "existinguser",
				Email:    "new@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			dbError: errors.New("Error 1062: Duplicate entry 'existinguser' for key 'username'"),
			wantErr: true,
		},
		{
			name: "Duplicate email",
			user: &model.User{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			dbError: errors.New("Error 1062: Duplicate entry 'existing@example.com' for key 'email'"),
			wantErr: true,
		},
		{
			name: "Missing required fields",
			user: &model.User{
				Username: "",
				Email:    "",
				Password: "",
			},
			dbError: errors.New("Error 1048: Column 'username' cannot be null"),
			wantErr: true,
		},
		{
			name: "Database connection error",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			dbError: errors.New("connection refused"),
			wantErr: true,
		},
		{
			name:    "Empty user object",
			user:    nil,
			dbError: errors.New("invalid user object: nil"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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

			store := &UserStore{
				db: gormDB,
			}

			if tt.user != nil {

				if tt.dbError != nil {
					mock.ExpectBegin()
					mock.ExpectExec("INSERT INTO `users`").
						WillReturnError(tt.dbError)
					mock.ExpectRollback()
				} else {
					mock.ExpectBegin()
					mock.ExpectExec("INSERT INTO `users`").
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				}
			}

			err = store.Create(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.dbError != nil {
					assert.Contains(t, err.Error(), tt.dbError.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed. Error: %v", tt.name, err)
		})
	}
}


/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1


 */
func TestUserStoreGetByID(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:   "Successfully retrieve user by valid ID",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(1, time.Now(), time.Now(), nil, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg")
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`deleted_at` IS NULL AND \\(`users`.`id` = \\?\\)").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
		{
			name:   "Attempt to retrieve non-existent user ID",
			userID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`deleted_at` IS NULL AND \\(`users`.`id` = \\?\\)").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:   "Handle database connection error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users`").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:   "Retrieve user with zero ID",
			userID: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`deleted_at` IS NULL AND \\(`users`.`id` = \\?\\)").
					WithArgs(0).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:   "Handle maximum ID value",
			userID: math.MaxUint32,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`deleted_at` IS NULL AND \\(`users`.`id` = \\?\\)").
					WithArgs(math.MaxUint32).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			tt.mockSetup(mock)

			store := &UserStore{db: gormDB}
			user, err := store.GetByID(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
				assert.Equal(t, tt.expectedUser.Image, user.Image)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1


 */
func TestUserStoreGetByEmail(t *testing.T) {

	type testCase struct {
		name          string
		email         string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}

	testTime := time.Now()
	validUser := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: testTime,
			UpdatedAt: testTime,
		},
		Email:    "test@example.com",
		Username: "testuser",
		Password: "hashedpassword",
		Bio:      "Test bio",
		Image:    "test-image.jpg",
	}

	tests := []testCase{
		{
			name:  "Successful user retrieval",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}).
					AddRow(validUser.ID, validUser.CreatedAt, validUser.UpdatedAt, nil, validUser.Username, validUser.Email, validUser.Password, validUser.Bio, validUser.Image)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (email = ?) ORDER BY "users"."id" ASC LIMIT 1`)).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedUser:  validUser,
			expectedError: nil,
		},
		{
			name:  "User not found",
			email: "nonexistent@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (email = ?) ORDER BY "users"."id" ASC LIMIT 1`)).
					WithArgs("nonexistent@example.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Empty email",
			email: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (email = ?) ORDER BY "users"."id" ASC LIMIT 1`)).
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (email = ?) ORDER BY "users"."id" ASC LIMIT 1`)).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			userStore := &UserStore{db: gormDB}

			user, err := userStore.GetByEmail(tc.email)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectedUser.Email, user.Email)
				assert.Equal(t, tc.expectedUser.Username, user.Username)
				assert.Equal(t, tc.expectedUser.Bio, user.Bio)
				assert.Equal(t, tc.expectedUser.Image, user.Image)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24


 */
func TestUserStoreGetByUsername(t *testing.T) {
	type testCase struct {
		name          string
		username      string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}

	now := time.Now()
	tests := []testCase{
		{
			name:     "Successfully retrieve user by valid username",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, now, now, nil, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg"))
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
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
			name:     "Handle non-existent username",
			username: "nonexistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("nonexistent").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Handle empty username",
			username: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Handle database connection error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("testuser").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlite3", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.mockSetup(mock)

			userStore := &UserStore{db: gormDB}

			user, err := userStore.GetByUsername(tc.username)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectedUser.ID, user.ID)
				assert.Equal(t, tc.expectedUser.Username, user.Username)
				assert.Equal(t, tc.expectedUser.Email, user.Email)
				assert.Equal(t, tc.expectedUser.Bio, user.Bio)
				assert.Equal(t, tc.expectedUser.Image, user.Image)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435


 */
func TestUserStoreUpdate(t *testing.T) {
	type testCase struct {
		name        string
		user        *model.User
		mockSetup   func(sqlmock.Sqlmock)
		expectedErr error
	}

	baseUser := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Bio:      "Test bio",
		Image:    "https://example.com/image.jpg",
	}

	tests := []testCase{
		{
			name: "Successful Update",
			user: baseUser,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `users` SET").
					WithArgs(
						sqlmock.AnyArg(),
						"testuser",
						"test@example.com",
						"password123",
						"Test bio",
						"https://example.com/image.jpg",
						1,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "Database Error",
			user: baseUser,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `users` SET").
					WithArgs(
						sqlmock.AnyArg(),
						"testuser",
						"test@example.com",
						"password123",
						"Test bio",
						"https://example.com/image.jpg",
						1,
					).WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectedErr: errors.New("database error"),
		},
		{
			name: "Invalid User Data",
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "",
				Email:    "",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `users` SET").
					WithArgs(
						sqlmock.AnyArg(),
						"",
						"",
						"",
						"",
						"",
						1,
					).WillReturnError(errors.New("validation error"))
				mock.ExpectRollback()
			},
			expectedErr: errors.New("validation error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			gormDB.LogMode(false)
			defer gormDB.Close()

			tc.mockSetup(mock)

			userStore := &UserStore{db: gormDB}

			err = userStore.Update(tc.user)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Follow_48fdf1257b
ROOST_METHOD_SIG_HASH=Follow_8217e61c06


 */
func TestUserStoreFollow(t *testing.T) {

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

	store := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		setupMock   func()
		follower    *model.User
		following   *model.User
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Follow",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			follower: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "user1@test.com",
			},
			following: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "user2",
				Email:    "user2@test.com",
			},
			expectError: false,
		},
		{
			name: "Self Follow Attempt",
			setupMock: func() {

			},
			follower: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "user1@test.com",
			},
			following: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "user1@test.com",
			},
			expectError: true,
			errorMsg:    "cannot follow self",
		},
		{
			name: "Follow with Nil User",
			setupMock: func() {

			},
			follower:    nil,
			following:   nil,
			expectError: true,
			errorMsg:    "invalid user objects",
		},
		{
			name: "Database Error",
			setupMock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			follower: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "user1",
				Email:    "user1@test.com",
			},
			following: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "user2",
				Email:    "user2@test.com",
			},
			expectError: true,
			errorMsg:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock()

			err := store.Follow(tt.follower, tt.following)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			t.Logf("Test '%s' completed. Error expected: %v, Got error: %v",
				tt.name, tt.expectError, err)
		})
	}
}


/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c


 */
func TestUserStoreIsFollowing(t *testing.T) {

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

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		mockSetup   func(sqlmock.Sqlmock)
		expected    bool
		expectedErr error
	}{
		{
			name:  "Valid Following Relationship",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:    true,
			expectedErr: nil,
		},
		{
			name:  "Valid Non-Following Relationship",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User A",
			userA:       nil,
			userB:       &model.User{Model: gorm.Model{ID: 2}},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User B",
			userA:       &model.User{Model: gorm.Model{ID: 1}},
			userB:       nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:  "Database Error",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
			},
			expected:    false,
			expectedErr: errors.New("database error"),
		},
		{
			name:        "Both Users Nil",
			userA:       nil,
			userB:       nil,
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:  "Self-Following Check",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:  "Zero ID User Check",
			userA: &model.User{Model: gorm.Model{ID: 0}},
			userB: &model.User{Model: gorm.Model{ID: 0}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `follows`").
					WithArgs(0, 0).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:    false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.mockSetup(mock)

			result, err := userStore.IsFollowing(tt.userA, tt.userB)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed. Result: %v, Error: %v", tt.name, result, err)
		})
	}
}


/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55


 */
func TestUserStoreUnfollow(t *testing.T) {

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

	userStore := &UserStore{db: gormDB}

	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Unfollow",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:  "Null User Parameters",
			userA: nil,
			userB: &model.User{},
			setupMock: func(mock sqlmock.Sqlmock) {

			},
			expectError: true,
			errorMsg:    "primary key can't be nil",
		},
		{
			name: "Database Error",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			expectError: true,
			errorMsg:    "database error",
		},
		{
			name: "Self-Unfollow Attempt",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.setupMock(mock)

			err := userStore.Unfollow(tt.userA, tt.userB)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7


 */
func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	type testCase struct {
		name          string
		userID        uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedIDs   []uint
		expectedError error
	}

	tests := []testCase{
		{
			name:   "Successfully retrieve following user IDs",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: nil,
		},
		{
			name:   "User with no followers",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name:   "Database connection error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnError(errors.New("database connection error"))
			},
			expectedIDs:   []uint{},
			expectedError: errors.New("database connection error"),
		},
		{
			name:   "Invalid user ID",
			userID: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(0).
					WillReturnError(errors.New("invalid user ID"))
			},
			expectedIDs:   []uint{},
			expectedError: errors.New("invalid user ID"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			tc.mockSetup(mock)

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			userStore := &UserStore{
				db: gormDB,
			}

			testUser := &model.User{
				Model: gorm.Model{
					ID:        tc.userID,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				Bio:      "test bio",
				Image:    "test.jpg",
			}

			ids, err := userStore.GetFollowingUserIDs(testUser)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedIDs, ids)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

