package store

import (
	"reflect"
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/DATA-DOG/go-sqlmock"
)






type mockDB struct {
	updateError     error
	lastUpdatedUser *model.User
}


/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore 

 */
func (l customLogger) Print(v ...interface{}) {}

func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		wantType reflect.Type
	}{
		{
			name:     "Valid gorm.DB instance",
			db:       &gorm.DB{},
			wantNil:  false,
			wantType: reflect.TypeOf(&UserStore{}),
		},
		{
			name:     "Nil gorm.DB instance",
			db:       nil,
			wantNil:  false,
			wantType: reflect.TypeOf(&UserStore{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewUserStore() returned nil: %v, want nil: %v", got == nil, tt.wantNil)
			}

			if got != nil && reflect.TypeOf(got) != tt.wantType {
				t.Errorf("NewUserStore() returned type %v, want %v", reflect.TypeOf(got), tt.wantType)
			}

			if got != nil && got.db != tt.db {
				t.Errorf("NewUserStore().db = %v, want %v", got.db, tt.db)
			}
		})
	}
}

func TestNewUserStoreImmutability(t *testing.T) {
	mockDB := &gorm.DB{}
	store1 := NewUserStore(mockDB)
	store2 := NewUserStore(mockDB)

	if store1 == store2 {
		t.Errorf("NewUserStore() returned the same instance for different calls")
	}
}

func TestNewUserStorePerformance(t *testing.T) {
	mockDB := &gorm.DB{}
	iterations := 10000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		NewUserStore(mockDB)
	}

	duration := time.Since(start)
	t.Logf("Time taken for %d iterations: %v", iterations, duration)

	if duration > time.Second {
		t.Errorf("NewUserStore() took too long: %v for %d iterations", duration, iterations)
	}
}

func TestNewUserStoreWithDifferentConfigurations(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
	}{
		{
			name: "Default configuration",
			db:   &gorm.DB{},
		},
		{
			name: "Custom logger",
			db:   &gorm.DB{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)

			if got == nil {
				t.Fatalf("NewUserStore() returned nil")
			}

			if got.db != tt.db {
				t.Errorf("NewUserStore().db = %v, want %v", got.db, tt.db)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1

FUNCTION_DEF=func (s *UserStore) GetByID(id uint) (*model.User, error) 

 */
func TestUserStoreGetById(t *testing.T) {
	tests := []struct {
		name     string
		id       uint
		wantUser *model.User
		wantErr  error
		setupDB  func(*gorm.DB)
	}{
		{
			name: "Successfully retrieve an existing user",
			id:   1,
			wantUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			wantErr: nil,
			setupDB: func(db *gorm.DB) {
				db.Create(&model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password",
				})
			},
		},
		{
			name:     "Attempt to retrieve a non-existent user",
			id:       999,
			wantUser: nil,
			wantErr:  gorm.ErrRecordNotFound,
			setupDB:  func(db *gorm.DB) {},
		},
		{
			name:     "Handle database connection error",
			id:       1,
			wantUser: nil,
			wantErr:  errors.New("database connection error"),
			setupDB: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
		},
		{
			name: "Retrieve a user with minimum fields populated",
			id:   2,
			wantUser: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "minuser",
				Email:    "min@example.com",
				Password: "minpass",
			},
			wantErr: nil,
			setupDB: func(db *gorm.DB) {
				db.Create(&model.User{
					Model:    gorm.Model{ID: 2},
					Username: "minuser",
					Email:    "min@example.com",
					Password: "minpass",
				})
			},
		},
		{
			name: "Retrieve a user with all fields populated",
			id:   3,
			wantUser: &model.User{
				Model:    gorm.Model{ID: 3},
				Username: "fulluser",
				Email:    "full@example.com",
				Password: "fullpass",
				Bio:      "Full bio",
				Image:    "full.jpg",
			},
			wantErr: nil,
			setupDB: func(db *gorm.DB) {
				db.Create(&model.User{
					Model:    gorm.Model{ID: 3},
					Username: "fulluser",
					Email:    "full@example.com",
					Password: "fullpass",
					Bio:      "Full bio",
					Image:    "full.jpg",
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, _ := gorm.Open("sqlite3", ":memory:")
			defer db.Close()
			db.AutoMigrate(&model.User{})
			tt.setupDB(db)

			s := &UserStore{db: db}

			gotUser, err := s.GetByID(tt.id)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("UserStore.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("UserStore.GetByID() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1

FUNCTION_DEF=func (s *UserStore) GetByEmail(email string) (*model.User, error) 

 */
func TestUserStoreGetByEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		mockDB  func() *gorm.DB
		want    *model.User
		wantErr error
	}{
		{
			name:  "Successfully retrieve a user by email",
			email: "user@example.com",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = nil
				return db.InstantSet("gorm:get_first_result", &model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
					Email:    "user@example.com",
				})
			},
			want: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "user@example.com",
			},
			wantErr: nil,
		},
		{
			name:  "Attempt to retrieve a non-existent user",
			email: "nonexistent@example.com",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = gorm.ErrRecordNotFound
				return db
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:  "Handle database connection error",
			email: "user@example.com",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = errors.New("database connection error")
				return db
			},
			want:    nil,
			wantErr: errors.New("database connection error"),
		},
		{
			name:  "Retrieve user with empty email string",
			email: "",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = gorm.ErrRecordNotFound
				return db
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:  "Case sensitivity in email lookup",
			email: "User@Example.com",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = nil
				return db.InstantSet("gorm:get_first_result", &model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
					Email:    "User@Example.com",
				})
			},
			want: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "User@Example.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserStore{
				db: tt.mockDB(),
			}
			got, err := s.GetByEmail(tt.email)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("UserStore.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("UserStore.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStore.GetByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24

FUNCTION_DEF=func (s *UserStore) GetByUsername(username string) (*model.User, error) 

 */
func TestUserStoreGetByUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		mockDB   func() *gorm.DB
		want     *model.User
		wantErr  error
	}{
		{
			name:     "Successfully retrieve a user by username",
			username: "testuser",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:get_struct", &model.User{Username: "testuser", Email: "test@example.com"})
			},
			want:    &model.User{Username: "testuser", Email: "test@example.com"},
			wantErr: nil,
		},
		{
			name:     "Attempt to retrieve a non-existent user",
			username: "nonexistent",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = gorm.ErrRecordNotFound
				return db
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "Handle database connection error",
			username: "testuser",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = errors.New("database connection error")
				return db
			},
			want:    nil,
			wantErr: errors.New("database connection error"),
		},
		{
			name:     "Retrieve user with maximum length username",
			username: string(make([]byte, 255)),
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:get_struct", &model.User{Username: string(make([]byte, 255)), Email: "max@example.com"})
			},
			want:    &model.User{Username: string(make([]byte, 255)), Email: "max@example.com"},
			wantErr: nil,
		},
		{
			name:     "Attempt to retrieve with an empty username",
			username: "",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = gorm.ErrRecordNotFound
				return db
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "Verify case sensitivity of username lookup",
			username: "TestUser",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:get_struct", &model.User{Username: "TestUser", Email: "case@example.com"})
			},
			want:    &model.User{Username: "TestUser", Email: "case@example.com"},
			wantErr: nil,
		},
		{
			name:     "Handle special characters in username",
			username: "test@user_123",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:get_struct", &model.User{Username: "test@user_123", Email: "special@example.com"})
			},
			want:    &model.User{Username: "test@user_123", Email: "special@example.com"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserStore{
				db: tt.mockDB(),
			}
			got, err := s.GetByUsername(tt.username)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("UserStore.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("UserStore.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStore.GetByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435

FUNCTION_DEF=func (s *UserStore) Update(m *model.User) error 

 */
func (m *mockDB) Error() error {
	return m.updateError
}

func (m *mockDB) Model(value interface{}) *mockDB {
	return m
}

func TestUserStoreUpdate(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		dbError error
		wantErr bool
	}{
		{
			name: "Successful Update",
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "updateduser",
				Email:    "updated@example.com",
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Update with No Changes",
			user: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "nochangeuser",
				Email:    "nochange@example.com",
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Update with Database Error",
			user: &model.User{
				Model:    gorm.Model{ID: 3},
				Username: "erroruser",
				Email:    "error@example.com",
			},
			dbError: errors.New("database error"),
			wantErr: true,
		},
		{
			name: "Update with Invalid User Data",
			user: &model.User{
				Model:    gorm.Model{ID: 4},
				Username: "",
				Email:    "invalid@example.com",
			},
			dbError: errors.New("invalid data"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &mockDB{
				updateError: tt.dbError,
			}

			store := &mockUserStore{
				db: mockDB,
			}

			err := store.Update(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(mockDB.lastUpdatedUser, tt.user) {
				t.Errorf("Update() called with user = %v, want %v", mockDB.lastUpdatedUser, tt.user)
			}
		})
	}
}

func (m *mockDB) Update(attrs ...interface{}) *mockDB {
	if len(attrs) > 0 {
		if user, ok := attrs[0].(*model.User); ok {
			m.lastUpdatedUser = user
		}
	}
	return m
}


/*
ROOST_METHOD_HASH=Follow_48fdf1257b
ROOST_METHOD_SIG_HASH=Follow_8217e61c06

FUNCTION_DEF=func (s *UserStore) Follow(a *model.User, b *model.User) error 

 */
func TestUserStoreFollow(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(*gorm.DB) (*model.User, *model.User, error)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Successful Follow Operation",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User, error) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				userB := &model.User{Username: "userB", Email: "userB@example.com"}
				db.Create(userA)
				db.Create(userB)
				return userA, userB, nil
			},
			expectError: false,
		},
		{
			name: "Follow a User That Is Already Followed",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User, error) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				userB := &model.User{Username: "userB", Email: "userB@example.com"}
				db.Create(userA)
				db.Create(userB)
				db.Model(userA).Association("Follows").Append(userB)
				return userA, userB, nil
			},
			expectError: false,
		},
		{
			name: "Attempt to Follow Non-Existent User",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User, error) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				db.Create(userA)
				userB := &model.User{Username: "userB", Email: "userB@example.com"}
				return userA, userB, nil
			},
			expectError: true,
			errorMsg:    "record not found",
		},
		{
			name: "Self-Follow Attempt",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User, error) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				db.Create(userA)
				return userA, userA, nil
			},
			expectError: false,
		},
		{
			name: "Follow Operation with Database Error",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User, error) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				userB := &model.User{Username: "userB", Email: "userB@example.com"}
				db.Create(userA)
				db.Create(userB)

				db.AddError(errors.New("database error"))
				return userA, userB, nil
			},
			expectError: true,
			errorMsg:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			require.NoError(t, err)
			defer db.Close()

			db.AutoMigrate(&model.User{})

			userStore := &UserStore{db: db}

			userA, userB, err := tt.setupFunc(db)
			require.NoError(t, err)

			err = userStore.Follow(userA, userB)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)

				var followedUsers []model.User
				db.Model(userA).Association("Follows").Find(&followedUsers)
				assert.Contains(t, followedUsers, *userB)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c

FUNCTION_DEF=func (s *UserStore) IsFollowing(a *model.User, b *model.User) (bool, error) 

 */
func TestUserStoreIsFollowing(t *testing.T) {
	type mockDB struct {
		count int
		err   error
	}

	tests := []struct {
		name     string
		userA    *model.User
		userB    *model.User
		mockDB   mockDB
		expected bool
		wantErr  bool
	}{
		{
			name:     "Following relationship exists",
			userA:    &model.User{Model: gorm.Model{ID: 1}},
			userB:    &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   mockDB{count: 1, err: nil},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Following relationship does not exist",
			userA:    &model.User{Model: gorm.Model{ID: 1}},
			userB:    &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   mockDB{count: 0, err: nil},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Database error",
			userA:    &model.User{Model: gorm.Model{ID: 1}},
			userB:    &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   mockDB{count: 0, err: errors.New("database error")},
			expected: false,
			wantErr:  true,
		},
		{
			name:     "User A is nil",
			userA:    nil,
			userB:    &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   mockDB{count: 0, err: nil},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "User B is nil",
			userA:    &model.User{Model: gorm.Model{ID: 1}},
			userB:    nil,
			mockDB:   mockDB{count: 0, err: nil},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Both users are nil",
			userA:    nil,
			userB:    nil,
			mockDB:   mockDB{count: 0, err: nil},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Self-following check",
			userA:    &model.User{Model: gorm.Model{ID: 1}},
			userB:    &model.User{Model: gorm.Model{ID: 1}},
			mockDB:   mockDB{count: 0, err: nil},
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &gorm.DB{}
			mockDB.AddError(tt.mockDB.err)

			us := &UserStore{db: mockDB}

			mockDB.Callback().Query().Register("mock_count", func(scope *gorm.Scope) {
				if tt.mockDB.err == nil {
					scope.InstanceSet("gorm:query_count", tt.mockDB.count)
				}
			})

			got, err := us.IsFollowing(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.IsFollowing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.expected {
				t.Errorf("UserStore.IsFollowing() = %v, want %v", got, tt.expected)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7

FUNCTION_DEF=func (s *UserStore) GetFollowingUserIDs(m *model.User) ([]uint, error) 

 */
func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		mockDB  func(mock sqlmock.Sqlmock)
		want    []uint
		wantErr bool
	}{
		{
			name: "Successful Retrieval of Following User IDs",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    []uint{2, 3, 4},
			wantErr: false,
		},
		{
			name: "User with No Followers",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    []uint{},
			wantErr: false,
		},
		{
			name: "Database Error Handling",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Large Number of Followers",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				for i := 2; i <= 10001; i++ {
					rows.AddRow(uint(i))
				}
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: func() []uint {
				ids := make([]uint, 10000)
				for i := range ids {
					ids[i] = uint(i + 2)
				}
				return ids
			}(),
			wantErr: false,
		},
		{
			name: "Invalid User ID",
			user: &model.User{Model: gorm.Model{ID: 0}},
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(0).
					WillReturnRows(rows)
			},
			want:    []uint{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gdb, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a gorm database", err)
			}
			defer gdb.Close()

			tt.mockDB(mock)

			s := &UserStore{db: gdb}

			got, err := s.GetFollowingUserIDs(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.GetFollowingUserIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStore.GetFollowingUserIDs() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

