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
	time "time"
)








/*
ROOST_METHOD_HASH=NewUserStore_fb599438e5
ROOST_METHOD_SIG_HASH=NewUserStore_c0075221af

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore // NewUserStore returns a new UserStore


*/
func TestNewUserStore(t *testing.T) {

	type testCase struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		validate func(*testing.T, *UserStore)
	}

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	gormDB, err := gorm.Open("mysql", mockDB)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	tests := []testCase{
		{
			name: "Successful UserStore Creation",
			db:   gormDB,
			validate: func(t *testing.T, us *UserStore) {
				if us == nil {
					t.Error("Expected non-nil UserStore")
					return
				}
				if us.db == nil {
					t.Error("Expected non-nil db in UserStore")
					return
				}
				t.Log("Successfully created UserStore with valid DB connection")
			},
		},
		{
			name:    "Nil DB Connection",
			db:      nil,
			wantNil: false,
			validate: func(t *testing.T, us *UserStore) {
				if us == nil {
					t.Error("Expected non-nil UserStore even with nil DB")
					return
				}
				if us.db != nil {
					t.Error("Expected nil db in UserStore when initialized with nil")
					return
				}
				t.Log("Successfully handled nil DB connection case")
			},
		},
		{
			name: "DB Reference Integrity",
			db:   gormDB,
			validate: func(t *testing.T, us *UserStore) {
				if us.db != gormDB {
					t.Error("DB reference doesn't match the provided instance")
					return
				}
				t.Log("Successfully verified DB reference integrity")
			},
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

			got := NewUserStore(tc.db)

			if (got == nil) != tc.wantNil {
				t.Errorf("NewUserStore() nil check = %v, want %v", got == nil, tc.wantNil)
				return
			}

			if tc.validate != nil {
				tc.validate(t, got)
			}
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
		wantErr bool
		setup   func(mock gosqlmock.Sqlmock)
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
			wantErr: false,
			setup: func(mock gosqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WithArgs(
						gosqlmock.AnyArg(),
						gosqlmock.AnyArg(),
						gosqlmock.AnyArg(),
						gosqlmock.AnyArg(),
						"testuser",
						"test@example.com",
						"password123",
						"Test bio",
						"https://example.com/image.jpg",
					).WillReturnResult(gosqlmock.NewResult(1, 1))
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
					WillReturnError(&sql.Error{
						Number:  1062,
						Message: "Duplicate entry 'existinguser' for key 'username'",
					})
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
					WillReturnError(&sql.Error{
						Number:  1062,
						Message: "Duplicate entry 'existing@example.com' for key 'email'",
					})
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
					WillReturnError(&sql.Error{
						Number:  1048,
						Message: "Column 'username' cannot be null",
					})
				mock.ExpectRollback()
			},
		},
		{
			name: "Database Connection Failure",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
			setup: func(mock gosqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnError(&sql.Error{
						Number:  1045,
						Message: "Connection refused",
					})
				mock.ExpectRollback()
			},
		},
		{
			name: "Large Data Fields",
			user: &model.User{
				Username: string(make([]byte, 255)),
				Email:    "test@example.com",
				Password: string(make([]byte, 255)),
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 1000)),
			},
			wantErr: false,
			setup: func(mock gosqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WithArgs(
						gosqlmock.AnyArg(),
						gosqlmock.AnyArg(),
						gosqlmock.AnyArg(),
						gosqlmock.AnyArg(),
						gosqlmock.AnyArg(),
						"test@example.com",
						gosqlmock.AnyArg(),
						gosqlmock.AnyArg(),
						gosqlmock.AnyArg(),
					).WillReturnResult(gosqlmock.NewResult(1, 1))
				mock.ExpectCommit()
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
				t.Logf("Test case '%s' failed with error: %v", tt.name, err)
			} else {
				t.Logf("Test case '%s' passed successfully", tt.name)
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
			name:  "Nil User A",
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
			name: "Nil User B",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: nil,
			setupMock: func(mock sqlmock.Sqlmock) {

			},
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
				t.Fatalf("Failed to open GORM DB: %v", err)
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

			t.Logf("Test '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_GetByEmail_fda09af5c4
ROOST_METHOD_SIG_HASH=UserStore_GetByEmail_9e84f3286b

FUNCTION_DEF=func (s *UserStore) GetByEmail(email string) (*model.User, error) // GetByEmail finds a user from email


*/
func TestUserStoreGetByEmail(t *testing.T) {

	tests := []struct {
		name          string
		email         string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}{
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

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			tt.mockSetup(mock)

			store := &UserStore{
				db: gormDB,
			}

			t.Logf("Testing scenario: %s", tt.name)
			user, err := store.GetByEmail(tt.email)

			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expectedUser != nil {
				if user == nil {
					t.Error("Expected user to not be nil")
				} else if user.Email != tt.expectedUser.Email {
					t.Errorf("Expected email %s, got %s", tt.expectedUser.Email, user.Email)
				}
			} else if user != nil {
				t.Error("Expected user to be nil")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_GetByID_1f5f06165b
ROOST_METHOD_SIG_HASH=UserStore_GetByID_2a864916bb

FUNCTION_DEF=func (s *UserStore) GetByID(id uint) (*model.User, error) // GetByID finds a user from id


*/
func TestUserStoreGetById(t *testing.T) {

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
		name    string
		id      uint
		mock    func()
		want    *model.User
		wantErr error
	}{
		{
			name: "Successfully retrieve existing user",
			id:   1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			want: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			wantErr: nil,
		},
		{
			name: "Attempt to retrieve non-existent user",
			id:   999,
			mock: func() {
				mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Handle database connection error",
			id:   1,
			mock: func() {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: sql.ErrConnDone,
		},
		{
			name: "Retrieve user with zero ID",
			id:   0,
			mock: func() {
				mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
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

			t.Log("Starting test case:", tt.name)

			tt.mock()

			got, err := userStore.GetByID(tt.id)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("UserStore.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("UserStore.GetByID() error message = %v, wantErr message %v", err, tt.wantErr)
				return
			}

			if got != nil && tt.want != nil {
				if got.ID != tt.want.ID ||
					got.Username != tt.want.Username ||
					got.Email != tt.want.Email ||
					got.Password != tt.want.Password ||
					got.Bio != tt.want.Bio ||
					got.Image != tt.want.Image {
					t.Errorf("UserStore.GetByID() = %v, want %v", got, tt.want)
				}
			} else if (got == nil) != (tt.want == nil) {
				t.Errorf("UserStore.GetByID() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Log("Successfully completed test case:", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_GetByUsername_622b1b9e41
ROOST_METHOD_SIG_HASH=UserStore_GetByUsername_992f00baec

FUNCTION_DEF=func (s *UserStore) GetByUsername(username string) (*model.User, error) // GetByUsername finds a user from username


*/
func TestUserStoreGetByUsername(t *testing.T) {

	tests := []struct {
		name     string
		username string
		mockDB   func(mock sqlmock.Sqlmock)
		want     *model.User
		wantErr  error
	}{
		{
			name:     "Success - Valid Username",
			username: "testuser",
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "bio", "image"}).
					AddRow(1, "testuser", "test@example.com", "hashedpass", "test bio", "image.jpg")
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			want: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			wantErr: nil,
		},
		{
			name:     "Error - User Not Found",
			username: "nonexistent",
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("nonexistent").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "Error - Empty Username",
			username: "",
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "Error - Database Connection Error",
			username: "testuser",
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WithArgs("testuser").
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: sql.ErrConnDone,
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

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer gormDB.Close()

			tt.mockDB(mock)

			store := &UserStore{
				db: gormDB,
			}

			got, err := store.GetByUsername(tt.username)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if got == nil && tt.want != nil {
					t.Error("GetByUsername() returned nil user, wanted non-nil")
					return
				}
				if got != nil && tt.want != nil {

					if got.Username != tt.want.Username {
						t.Errorf("GetByUsername() username = %v, want %v", got.Username, tt.want.Username)
					}
					if got.Email != tt.want.Email {
						t.Errorf("GetByUsername() email = %v, want %v", got.Email, tt.want.Email)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
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
			name:   "Success - Multiple Following Users",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)
				mock.ExpectQuery("SELECT to_user_id FROM follows").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    []uint{2, 3, 4},
			wantErr: false,
		},
		{
			name:   "Success - Empty Following List",
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
			name:   "Error - Database Error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    []uint{},
			wantErr: true,
		},
		{
			name:   "Success - Large Dataset",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				for i := uint(2); i <= 1001; i++ {
					rows.AddRow(i)
				}
				mock.ExpectQuery("SELECT to_user_id FROM follows").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    generateSequentialIDs(2, 1001),
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

			if !compareUintSlices(got, tt.want) {
				t.Errorf("GetFollowingUserIDs() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test '%s' completed successfully", tt.name)
		})
	}
}

func compareUintSlices(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func generateSequentialIDs(start, end uint) []uint {
	var ids []uint
	for i := start; i <= end; i++ {
		ids = append(ids, i)
	}
	return ids
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

			t.Logf("Running test case: %s", tt.name)

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

	type testCase struct {
		name        string
		userA       *model.User
		userB       *model.User
		setupMock   func(mock sqlmock.Sqlmock)
		expectError bool
	}

	tests := []testCase{
		{
			name: "Successful Unfollow",
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
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Unfollow Non-Followed User",
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
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Invalid User A",
			userA: &model.User{
				Model: gorm.Model{ID: 0},
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(0, 2).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
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
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnError(gorm.ErrInvalidTransaction)
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
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer gormDB.Close()

			tc.setupMock(mock)

			store := &UserStore{db: gormDB}

			err = store.Unfollow(tc.userA, tc.userB)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tc.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_Update_4fd6d3d1c1
ROOST_METHOD_SIG_HASH=UserStore_Update_ddd5c151cf

FUNCTION_DEF=func (s *UserStore) Update(m *model.User) error // Update update all of user fields


*/
func TestUserStoreUpdate(t *testing.T) {

	tests := []struct {
		name      string
		user      *model.User
		setupMock func(sqlmock.Sqlmock)
		wantErr   bool
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
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WithArgs(
						sqlmock.AnyArg(),
						"testuser",
						"test@example.com",
						"password123",
						"Test bio",
						"test.jpg",
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
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Database Connection Error",
			user: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Empty User Model",
			user: &model.User{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(gorm.ErrRecordNotFound)
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

			tt.setupMock(mock)

			store := &UserStore{
				db: gormDB,
			}

			err = store.Update(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if !tt.wantErr {
				t.Logf("Successfully updated user with ID: %d", tt.user.ID)
			}
		})
	}
}

