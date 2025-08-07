package store

import (
	debug "runtime/debug"
	testing "testing"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	sql "database/sql"
)








/*
ROOST_METHOD_HASH=NewUserStore_fb599438e5
ROOST_METHOD_SIG_HASH=NewUserStore_c0075221af

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore // NewUserStore returns a new UserStore


*/
func TestNewUserStore(t *testing.T) {

	tests := []struct {
		name    string
		db      *gorm.DB
		wantErr bool
	}{
		{
			name:    "Successfully create new UserStore instance",
			db:      nil,
			wantErr: false,
		},
		{
			name:    "Handle nil database connection",
			db:      nil,
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

			t.Log("Starting test case:", tt.name)

			if tt.db == nil && !tt.wantErr {

				mockDB, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("Failed to create mock DB: %v", err)
				}
				defer mockDB.Close()

				gormDB, err := gorm.Open("mysql", mockDB)
				if err != nil {
					t.Fatalf("Failed to create GORM DB: %v", err)
				}
				tt.db = gormDB

				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("Mock expectations were not met: %v", err)
				}
			}

			got := NewUserStore(tt.db)

			if !tt.wantErr {
				if got == nil {
					t.Error("NewUserStore() returned nil, expected valid UserStore instance")
				}
				if got.db != tt.db {
					t.Error("NewUserStore() database connection mismatch")
				}
				t.Log("Successfully created UserStore instance")
			} else {
				if got != nil && got.db != nil {
					t.Error("NewUserStore() expected to handle nil database gracefully")
				}
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

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{
		db: gormDB,
	}

	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
		setup   func(sqlmock.Sqlmock)
	}{
		{
			name: "Successful user creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs(
						sqlmock.AnyArg(),
						"testuser",
						"test@example.com",
						"password123",
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).WillReturnResult(sqlmock.NewResult(1, 1))
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

			t.Logf("Running test case: %s", tt.name)

			tt.setup(mock)

			err := userStore.Create(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case completed successfully: %s", tt.name)
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
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{
		db: gormDB,
	}

	tests := []struct {
		name    string
		userID  uint
		want    *model.User
		wantErr bool
	}{
		{
			name:   "Successfully retrieve existing user",
			userID: 1,
			want: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
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

			t.Logf("Testing scenario: %s", tt.name)

			if !tt.wantErr {
				rows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(tt.want.ID, tt.want.Username, tt.want.Email)
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			}

			got, err := userStore.GetByID(tt.userID)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Error("UserStore.GetByID() returned nil user when expecting success")
					return
				}

				if got.ID != tt.want.ID {
					t.Errorf("UserStore.GetByID() got ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Username != tt.want.Username {
					t.Errorf("UserStore.GetByID() got Username = %v, want %v", got.Username, tt.want.Username)
				}
				if got.Email != tt.want.Email {
					t.Errorf("UserStore.GetByID() got Email = %v, want %v", got.Email, tt.want.Email)
				}

				t.Logf("Successfully retrieved user with ID: %d", tt.userID)
			}
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
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		email       string
		mockSetup   func(sqlmock.Sqlmock)
		expectUser  *model.User
		expectError bool
	}{
		{
			name:  "Successfully retrieve user by valid email",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "username"}).
					AddRow(1, "test@example.com", "testuser")
				mock.ExpectQuery("SELECT").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectUser: &model.User{
				ID:       1,
				Email:    "test@example.com",
				Username: "testuser",
			},
			expectError: false,
		},
		{
			name:  "User not found",
			email: "nonexistent@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WithArgs("nonexistent@example.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectUser:  nil,
			expectError: true,
		},
		{
			name:  "Database error",
			email: "test@example.com",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WithArgs("test@example.com").
					WillReturnError(gorm.ErrInvalidSQL)
			},
			expectUser:  nil,
			expectError: true,
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

			tt.mockSetup(mock)

			t.Logf("Testing scenario: %s", tt.name)
			user, err := userStore.GetByEmail(tt.email)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.expectUser != nil {
				if user == nil {
					t.Error("Expected user but got nil")
				} else {
					if user.Email != tt.expectUser.Email {
						t.Errorf("Expected email %s but got %s", tt.expectUser.Email, user.Email)
					}
					if user.Username != tt.expectUser.Username {
						t.Errorf("Expected username %s but got %s", tt.expectUser.Username, user.Username)
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
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
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to create GORM DB instance: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{
		db: gormDB,
	}

	tests := []struct {
		name      string
		username  string
		mockSetup func(sqlmock.Sqlmock)
		wantUser  *model.User
		wantErr   bool
	}{
		{
			name:     "Successfully retrieve existing user",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "testuser", "test@example.com")
				mock.ExpectQuery("SELECT").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			wantUser: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WithArgs("nonexistent").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantUser: nil,
			wantErr:  true,
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

			tt.mockSetup(mock)

			gotUser, err := userStore.GetByUsername(tt.username)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if gotUser != nil {
					t.Errorf("UserStore.GetByUsername() = %v, want nil when error", gotUser)
				}
				t.Logf("Successfully verified error case")
				return
			}

			if gotUser == nil {
				t.Error("UserStore.GetByUsername() returned nil user, expected non-nil")
				return
			}

			if gotUser.Username != tt.wantUser.Username {
				t.Errorf("UserStore.GetByUsername() username = %v, want %v",
					gotUser.Username, tt.wantUser.Username)
			}

			if gotUser.Email != tt.wantUser.Email {
				t.Errorf("UserStore.GetByUsername() email = %v, want %v",
					gotUser.Email, tt.wantUser.Email)
			}

			t.Logf("Successfully verified user retrieval")
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
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
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		user        *model.User
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successful Update",
			user: &model.User{
				ID:       1,
				Username: "updated_user",
				Email:    "updated@example.com",
			},
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
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

			tt.setupMock(mock)

			err := userStore.Update(tt.user)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
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
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{
		db: gormDB,
	}

	tests := []struct {
		name    string
		userA   *model.User
		userB   *model.User
		wantErr bool
	}{
		{
			name: "Successful Follow Relationship",
			userA: &model.User{
				ID:       1,
				Username: "userA",
				Email:    "userA@example.com",
			},
			userB: &model.User{
				ID:       2,
				Username: "userB",
				Email:    "userB@example.com",
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

			t.Logf("Testing scenario: %s", tt.name)

			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			err := userStore.Follow(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Follow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err == nil {
				t.Log("Successfully created follow relationship")
			}
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
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		mockSetup   func(sqlmock.Sqlmock)
		expected    bool
		expectError bool
	}{
		{
			name: "Valid Following Relationship",
			userA: &model.User{
				ID: 1,
			},
			userB: &model.User{
				ID: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `follows` WHERE \\(from_user_id = \\? AND to_user_id = \\?\\)").
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:    true,
			expectError: false,
		},
		{
			name:  "Nil User A",
			userA: nil,
			userB: &model.User{
				ID: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {

			},
			expected:    false,
			expectError: false,
		},
		{
			name: "Database Error",
			userA: &model.User{
				ID: 1,
			},
			userB: &model.User{
				ID: 2,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `follows` WHERE \\(from_user_id = \\? AND to_user_id = \\?\\)").
					WithArgs(1, 2).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expected:    false,
			expectError: true,
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

			tt.mockSetup(mock)

			result, err := userStore.IsFollowing(tt.userA, tt.userB)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected result %v but got %v", tt.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=UserStore_Unfollow_29d3ef7f50
ROOST_METHOD_SIG_HASH=UserStore_Unfollow_31d9214353

FUNCTION_DEF=func (s *UserStore) Unfollow(a *model.User, b *model.User) error // Unfollow delete follow relashionship to User B from user A


*/
func TestUserStoreUnfollow(t *testing.T) {

	db, mock, err := gosqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	userStore := &UserStore{
		db: gormDB,
	}

	tests := []struct {
		name    string
		userA   *model.User
		userB   *model.User
		wantErr bool
		setup   func()
	}{
		{
			name: "Successful unfollow",
			userA: &model.User{
				ID:       1,
				Username: "userA",
				Email:    "userA@test.com",
			},
			userB: &model.User{
				ID:       2,
				Username: "userB",
				Email:    "userB@test.com",
			},
			wantErr: false,
			setup: func() {

				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(1, 2).
					WillReturnResult(gosqlmock.NewResult(0, 1))
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

			t.Log("Starting test case:", tt.name)

			if tt.setup != nil {
				tt.setup()
			}

			err := userStore.Unfollow(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Unfollow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Log("Successfully completed test case:", tt.name)
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
		userID        uint
		mockSetup     func(mock gosqlmock.Sqlmock)
		expectedIDs   []uint
		expectedError error
	}

	tests := []testCase{
		{
			name:   "Successfully retrieve multiple following user IDs",
			userID: 1,
			mockSetup: func(mock gosqlmock.Sqlmock) {
				rows := gosqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)

				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: nil,
		},
		{
			name:   "Empty following list",
			userID: 1,
			mockSetup: func(mock gosqlmock.Sqlmock) {
				rows := gosqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name:   "Database error",
			userID: 1,
			mockSetup: func(mock gosqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedIDs:   []uint{},
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

			db, mock, err := gosqlmock.New()
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

			userStore := &UserStore{
				db: gormDB,
			}

			testUser := &model.User{
				ID: tc.userID,
			}

			gotIDs, err := userStore.GetFollowingUserIDs(testUser)

			if err != tc.expectedError {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}

			if len(gotIDs) != len(tc.expectedIDs) {
				t.Errorf("Expected %d following IDs, got %d", len(tc.expectedIDs), len(gotIDs))
			}

			for i, id := range gotIDs {
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

