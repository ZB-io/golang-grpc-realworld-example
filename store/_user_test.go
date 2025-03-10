package store

import (
	debug "runtime/debug"
	testing "testing"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	assert "github.com/stretchr/testify/assert"
	sql "database/sql"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	time "time"
	errors "errors"
	fmt "fmt"
	sync "sync"
)








/*
ROOST_METHOD_HASH=NewUserStore_fb599438e5
ROOST_METHOD_SIG_HASH=NewUserStore_c0075221af

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore // NewUserStore returns a new UserStore


*/
func TestNewUserStore(t *testing.T) {

	tests := []struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}{
		{
			name: "Scenario 1: Successfully Create New UserStore Instance",
			db: func() *gorm.DB {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Failed to create mock DB: %v", err)
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("Failed to create GORM DB: %v", err)
				}

				mock.ExpectPing()
				return gormDB
			}(),
			wantNil:  false,
			scenario: "Basic constructor functionality validation",
		},
		{
			name:     "Scenario 2: Create UserStore with Nil Database Connection",
			db:       nil,
			wantNil:  false,
			scenario: "Edge case: nil database connection",
		},
		{
			name: "Scenario 3: Verify UserStore Instance Independence",
			db: func() *gorm.DB {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Failed to create mock DB: %v", err)
				}
				gormDB, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("Failed to create GORM DB: %v", err)
				}
				mock.ExpectPing()
				return gormDB
			}(),
			wantNil:  false,
			scenario: "Instance independence verification",
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

			t.Logf("Testing scenario: %s", tt.scenario)

			got := NewUserStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, got, "Expected nil UserStore")
			} else {
				assert.NotNil(t, got, "Expected non-nil UserStore")
				assert.Equal(t, tt.db, got.db, "Database connection mismatch")
			}

			switch tt.name {
			case "Scenario 3: Verify UserStore Instance Independence":

				got2 := NewUserStore(tt.db)
				assert.NotEqual(t, &got, &got2, "UserStore instances should have different memory addresses")
				assert.Equal(t, got.db, got2.db, "Database connections should be the same")
			}

			t.Logf("Successfully completed test: %s", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_Create_9495ddb29d
ROOST_METHOD_SIG_HASH=UserStore_Create_18451817fe

FUNCTION_DEF=func (s *UserStore) Create(m *model.User) error // Create create a user


*/
func TestUserStoreCreate(t *testing.T) {

	type testCase struct {
		name    string
		user    *model.User
		dbSetup func(mock sqlmock.Sqlmock)
		wantErr bool
	}

	tests := []testCase{
		{
			name: "Successful User Creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
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
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(&sql.Error{Message: "Duplicate entry 'existinguser' for key 'username'"})
				mock.ExpectRollback()
			},
			wantErr: true,
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
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(&sql.Error{Message: "Duplicate entry 'existing@example.com' for key 'email'"})
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Missing Required Fields",
			user: &model.User{
				Username: "",
				Email:    "",
				Password: "",
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(&sql.Error{Message: "Column 'username' cannot be null"})
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Database Connection Error",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(&sql.Error{Message: "connection refused"})
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Large Data Creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 1000)),
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
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

			tt.dbSetup(mock)

			store := &UserStore{db: gormDB}

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

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	store := &UserStore{db: gormDB}

	tests := []struct {
		name    string
		userA   *model.User
		userB   *model.User
		mockSQL func()
		wantErr bool
	}{
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
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
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
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:  "Invalid User A",
			userA: nil,
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			mockSQL: func() {

			},
			wantErr: true,
		},
		{
			name: "Invalid User B",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: nil,
			mockSQL: func() {

			},
			wantErr: true,
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
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
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

			t.Logf("Running test case: %s", tt.name)

			if tt.mockSQL != nil {
				tt.mockSQL()
			}

			err := store.Follow(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Follow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case completed: %s", tt.name)
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
			name:  "Success - Valid Email",
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
			name:  "Failure - Email Not Found",
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
			name:  "Failure - Database Error",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("test@example.com").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:  "Failure - Empty Email",
			email: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("").
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

			if err != tc.expectedError {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}

			if tc.expectedUser != nil {
				if user == nil {
					t.Error("Expected user to not be nil")
				} else if user.Email != tc.expectedUser.Email {
					t.Errorf("Expected email %s, got %s", tc.expectedUser.Email, user.Email)
				}
			} else if user != nil {
				t.Error("Expected user to be nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
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
		userID        uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}

	tests := []testCase{
		{
			name:   "Successfully retrieve existing user",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				columns := []string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "bio", "image"}
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, time.Now(), time.Now(), nil, "testuser", "test@example.com", "hashedpass", "bio", "image.jpg"))
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
			name:   "User not found",
			userID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:   "Database error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:   "Zero ID provided",
			userID: 0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs(0).
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
			user, err := store.GetByID(tc.userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, user)
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectedUser.ID, user.ID)
				assert.Equal(t, tc.expectedUser.Username, user.Username)
				assert.Equal(t, tc.expectedUser.Email, user.Email)
				t.Logf("Successfully retrieved user with ID: %d", user.ID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
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

	type testCase struct {
		name          string
		username      string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
		shouldPanic   bool
	}

	tests := []testCase{
		{
			name:     "Success - Valid Username",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "hashedpass", "bio", "image.jpg")
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpass",
				Bio:      "bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
		},
		{
			name:     "Failure - User Not Found",
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
			name:     "Failure - Empty Username",
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
			name:     "Failure - Database Error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("testuser").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:     "Success - Username with Special Characters",
			username: "test@user#123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "test@user#123", "special@example.com", "hashedpass", "bio", "image.jpg")
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("test@user#123").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Username: "test@user#123",
				Email:    "special@example.com",
				Password: "hashedpass",
				Bio:      "bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
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

			store := &UserStore{
				db: gormDB,
			}

			user, err := store.GetByUsername(tc.username)

			if err != tc.expectedError {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}

			if tc.expectedUser != nil {
				if user == nil {
					t.Error("Expected user to not be nil")
				} else {
					if user.Username != tc.expectedUser.Username {
						t.Errorf("Expected username %s, got %s", tc.expectedUser.Username, user.Username)
					}
					if user.Email != tc.expectedUser.Email {
						t.Errorf("Expected email %s, got %s", tc.expectedUser.Email, user.Email)
					}
				}
			} else if user != nil {
				t.Error("Expected user to be nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
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

	tests := []struct {
		name      string
		userID    uint
		mockSetup func(sqlmock.Sqlmock)
		want      []uint
		wantErr   bool
	}{
		{
			name:   "Successfully retrieve following user IDs",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(uint(2)).
					AddRow(uint(3)).
					AddRow(uint(4))
				mock.ExpectQuery("SELECT to_user_id FROM follows").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    []uint{2, 3, 4},
			wantErr: false,
		},
		{
			name:   "Empty following list",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    []uint{},
			wantErr: false,
		},
		{
			name:   "Database error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    []uint{},
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
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tt.mockSetup(mock)

			store := &UserStore{
				db: gormDB,
			}

			user := &model.User{
				Model: gorm.Model{ID: tt.userID},
			}

			got, err := store.GetFollowingUserIDs(user)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFollowingUserIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("GetFollowingUserIDs() got = %v, want %v", got, tt.want)
					return
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("GetFollowingUserIDs() got[%d] = %v, want[%d] = %v", i, got[i], i, tt.want[i])
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
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

	store := &UserStore{db: gormDB}

	tests := []struct {
		name     string
		userA    *model.User
		userB    *model.User
		mockFunc func()
		want     bool
		wantErr  error
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
			wantErr: nil,
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
			wantErr: nil,
		},
		{
			name:     "Nil User A",
			userA:    nil,
			userB:    &model.User{Model: gorm.Model{ID: 2}},
			mockFunc: func() {},
			want:     false,
			wantErr:  nil,
		},
		{
			name:     "Nil User B",
			userA:    &model.User{Model: gorm.Model{ID: 1}},
			userB:    nil,
			mockFunc: func() {},
			want:     false,
			wantErr:  nil,
		},
		{
			name:     "Both Users Nil",
			userA:    nil,
			userB:    nil,
			mockFunc: func() {},
			want:     false,
			wantErr:  nil,
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
			wantErr: sql.ErrConnDone,
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
			wantErr: nil,
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

			t.Logf("Running test case: %s", tt.name)

			tt.mockFunc()

			got, err := store.IsFollowing(tt.userA, tt.userB)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("IsFollowing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("IsFollowing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("IsFollowing() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case completed successfully: %s", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_Unfollow_29d3ef7f50
ROOST_METHOD_SIG_HASH=UserStore_Unfollow_31d9214353

FUNCTION_DEF=func (s *UserStore) Unfollow(a *model.User, b *model.User) error // Unfollow delete follow relashionship to User B from user A


*/
func TestUserStoreUnfollow(t *testing.T) {

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

	store := &UserStore{db: gormDB}

	tests := []struct {
		name    string
		userA   *model.User
		userB   *model.User
		mockSQL func()
		wantErr bool
	}{
		{
			name: "Successful unfollow",
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
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:  "Invalid user A (nil)",
			userA: nil,
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			mockSQL: func() {

			},
			wantErr: true,
		},
		{
			name: "Invalid user B (nil)",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB:   nil,
			mockSQL: func() {},
			wantErr: true,
		},
		{
			name: "Database error",
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
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(fmt.Errorf("database error"))
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

			t.Logf("Running test case: %s", tt.name)

			tt.mockSQL()

			err := store.Unfollow(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Unfollow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case completed: %s", tt.name)
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
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &UserStore{db: gormDB}

	tests := []struct {
		name    string
		user    *model.User
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successful Update",
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "updated_user",
				Email:    "updated@example.com",
				Bio:      "Updated bio",
				Image:    "updated.jpg",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WithArgs(
						sqlmock.AnyArg(),
						"updated_user",
						"updated@example.com",
						"Updated bio",
						"updated.jpg",
						1,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Non-Existent User",
			user: &model.User{
				Model:    gorm.Model{ID: 999},
				Username: "nonexistent",
				Email:    "nonexistent@example.com",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Duplicate Username",
			user: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "existing_user",
				Email:    "user2@example.com",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Empty Required Fields",
			user: &model.User{
				Model:    gorm.Model{ID: 3},
				Username: "",
				Email:    "",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(gorm.ErrInvalidSQL)
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

			tt.mockFn(mock)

			err := store.Update(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}

	t.Run("Concurrent Updates", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
				t.Fail()
			}
		}()

		user := &model.User{
			Model:    gorm.Model{ID: 4},
			Username: "concurrent_user",
			Email:    "concurrent@example.com",
		}

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `users`").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		var wg sync.WaitGroup
		updateCount := 3
		errChan := make(chan error, updateCount)

		for i := 0; i < updateCount; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				user.Bio = "Updated bio " + string(index)
				errChan <- store.Update(user)
			}(i)
		}

		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				t.Errorf("Concurrent update error: %v", err)
			}
		}

		t.Log("Concurrent update test completed")
	})
}

