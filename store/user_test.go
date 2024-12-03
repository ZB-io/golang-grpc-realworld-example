package store

import (
	"errors"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
	"sync"
	"database/sql"
)

/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920


 */
func TestCreateUser(t *testing.T) {
	type testCase struct {
		name    string
		user    *model.User
		dbError error
		wantErr bool
	}

	now := time.Now()

	tests := []testCase{
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
			name: "Duplicate username error",
			user: &model.User{
				Username: "existing_user",
				Email:    "another@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			dbError: errors.New("Error 1062: Duplicate entry"),
			wantErr: true,
		},
		{
			name: "Missing required fields",
			user: &model.User{
				Username: "",
				Email:    "",
				Password: "password123",
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
			name: "Invalid email format",
			user: &model.User{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			dbError: errors.New("invalid email format"),
			wantErr: true,
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
			defer gormDB.Close()

			store := &UserStore{
				db: gormDB,
			}

			if tc.dbError == nil {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						nil,
						tc.user.Username,
						tc.user.Email,
						tc.user.Password,
						tc.user.Bio,
						tc.user.Image,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			} else {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(tc.dbError)
				mock.ExpectRollback()
			}

			err = store.Create(tc.user)

			if tc.wantErr {
				assert.Error(t, err)
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
func (s *UserStore) Follow(a *model.User, b *model.User) error {
	if a == nil || b == nil {
		return errors.New("invalid user parameters")
	}

	if a.ID == b.ID {
		return errors.New("self-follow not allowed")
	}

	return s.db.Model(a).Association("Follows").Append(b).Error
}

func TestFollow(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	store := &UserStore{db: db}

	tests := []struct {
		name    string
		setup   func() (*model.User, *model.User)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful follow",
			setup: func() (*model.User, *model.User) {
				userA := &model.User{
					Username: "userA",
					Email:    "userA@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				userB := &model.User{
					Username: "userB",
					Email:    "userB@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				db.Create(userA)
				db.Create(userB)
				return userA, userB
			},
			wantErr: false,
		},
		{
			name: "Follow with nil user A",
			setup: func() (*model.User, *model.User) {
				userB := &model.User{
					Username: "userB2",
					Email:    "userB2@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				db.Create(userB)
				return nil, userB
			},
			wantErr: true,
			errMsg:  "invalid user parameters",
		},
		{
			name: "Follow with nil user B",
			setup: func() (*model.User, *model.User) {
				userA := &model.User{
					Username: "userA2",
					Email:    "userA2@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				db.Create(userA)
				return userA, nil
			},
			wantErr: true,
			errMsg:  "invalid user parameters",
		},
		{
			name: "Self follow attempt",
			setup: func() (*model.User, *model.User) {
				user := &model.User{
					Username: "selfUser",
					Email:    "self@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				db.Create(user)
				return user, user
			},
			wantErr: true,
			errMsg:  "self-follow not allowed",
		},
		{
			name: "Follow already followed user",
			setup: func() (*model.User, *model.User) {
				userA := &model.User{
					Username: "userA3",
					Email:    "userA3@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				userB := &model.User{
					Username: "userB3",
					Email:    "userB3@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				db.Create(userA)
				db.Create(userB)
				store.Follow(userA, userB)
				return userA, userB
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userA, userB := tt.setup()

			err := store.Follow(userA, userB)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)

				if userA != nil && userB != nil {
					var follows []model.User
					err = db.Model(userA).Association("Follows").Find(&follows).Error
					assert.NoError(t, err)

					found := false
					for _, f := range follows {
						if f.ID == userB.ID {
							found = true
							break
						}
					}
					assert.True(t, found, "Follow relationship not found in database")
				}
			}

			if userA != nil {
				db.Unscoped().Delete(userA)
			}
			if userB != nil {
				db.Unscoped().Delete(userB)
			}
		})
	}
}

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.User{})

	return db, nil
}

/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1


 */
func (m *mockDBForUserTest) First(out interface{}, where ...interface{}) *gorm.DB {
	called := m.Called(out, where)
	return called.Get(0).(*gorm.DB)
}

func TestGetByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		mockSetup     func(*mockDBForUserTest)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:  "Successfully retrieve user by valid email",
			email: "test@example.com",
			mockSetup: func(mock *mockDBForUserTest) {
				expectedUser := &model.User{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Email:    "test@example.com",
					Username: "testuser",
				}
				mock.On("Where", "email = ?", []interface{}{"test@example.com"}).
					Return(&gorm.DB{Error: nil})
				mock.On("First", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*model.User)
						*arg = *expectedUser
					}).
					Return(&gorm.DB{Error: nil})
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Email:    "test@example.com",
				Username: "testuser",
			},
			expectedError: nil,
		},
		{
			name:  "Non-existent email",
			email: "nonexistent@example.com",
			mockSetup: func(mock *mockDBForUserTest) {
				mock.On("Where", "email = ?", []interface{}{"nonexistent@example.com"}).
					Return(&gorm.DB{Error: nil})
				mock.On("First", mock.Anything, mock.Anything).
					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Empty email",
			email: "",
			mockSetup: func(mock *mockDBForUserTest) {
				mock.On("Where", "email = ?", []interface{}{""}).
					Return(&gorm.DB{Error: nil})
				mock.On("First", mock.Anything, mock.Anything).
					Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Database connection error",
			email: "test@example.com",
			mockSetup: func(mock *mockDBForUserTest) {
				mock.On("Where", "email = ?", []interface{}{"test@example.com"}).
					Return(&gorm.DB{Error: errors.New("database connection error")})
				mock.On("First", mock.Anything, mock.Anything).
					Return(&gorm.DB{Error: errors.New("database connection error")})
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
		{
			name:  "Special characters in email",
			email: "test+special@example.com",
			mockSetup: func(mock *mockDBForUserTest) {
				expectedUser := &model.User{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Email:    "test+special@example.com",
					Username: "testuser",
				}
				mock.On("Where", "email = ?", []interface{}{"test+special@example.com"}).
					Return(&gorm.DB{Error: nil})
				mock.On("First", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*model.User)
						*arg = *expectedUser
					}).
					Return(&gorm.DB{Error: nil})
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Email:    "test+special@example.com",
				Username: "testuser",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDBForUserTest)
			tt.mockSetup(mockDB)

			store := &UserStore{
				db: mockDB,
			}

			user, err := store.GetByEmail(tt.email)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *mockDBForUserTest) Where(query interface{}, args ...interface{}) *gorm.DB {
	called := m.Called(query, args)
	return called.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1


 */
func (m *mockDBForUser) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func TestGetByIDUser(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		setupMock     func(*mockDBForUser) *gorm.DB
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "Scenario 1: Successfully Retrieve User by Valid ID",
			id:   1,
			setupMock: func(m *mockDBForUser) *gorm.DB {
				expectedUser := &model.User{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Username: "testuser",
					Email:    "test@example.com",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				return &gorm.DB{Error: nil, Value: expectedUser}
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "testuser",
				Email:    "test@example.com",
				Bio:      "test bio",
				Image:    "test.jpg",
			},
			expectedError: nil,
		},
		{
			name: "Scenario 2: Non-existent User ID",
			id:   999,
			setupMock: func(m *mockDBForUser) *gorm.DB {
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Scenario 3: Database Connection Error",
			id:   1,
			setupMock: func(m *mockDBForUser) *gorm.DB {
				return &gorm.DB{Error: errors.New("database connection error")}
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
		{
			name: "Scenario 4: Zero ID Value",
			id:   0,
			setupMock: func(m *mockDBForUser) *gorm.DB {
				return &gorm.DB{Error: gorm.ErrRecordNotFound}
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Scenario 5: Soft-Deleted User",
			id:   2,
			setupMock: func(m *mockDBForUser) *gorm.DB {
				deletedAt := time.Now()
				return &gorm.DB{
					Error: gorm.ErrRecordNotFound,
					Value: &model.User{
						Model: gorm.Model{
							ID:        2,
							DeletedAt: &deletedAt,
						},
					},
				}
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDBForUser)
			store := &UserStore{db: &gorm.DB{}}

			db := tt.setupMock(mockDB)
			mockDB.On("Find", mock.Anything, mock.Anything).Return(db)

			user, err := store.GetByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
			}

			mockDB.AssertExpectations(t)
			t.Logf("Test case '%s' completed", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24


 */
func (m *mockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	if m.user != nil {
		if u, ok := out.(*model.User); ok {
			*u = *m.user
		}
	}
	return &gorm.DB{
		Error: m.error,
		Value: m.user,
	}
}

func TestGetByUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		mockDB   func() *mockDB
		want     *model.User
		wantErr  error
	}{
		{
			name:     "Successfully retrieve existing user",
			username: "testuser",
			mockDB: func() *mockDB {
				return &mockDB{
					user: &model.User{
						Model: gorm.Model{
							ID:        1,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Username: "testuser",
						Email:    "test@example.com",
						Bio:      "Test bio",
						Image:    "test-image.jpg",
					},
				}
			},
			want: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			wantErr: nil,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			mockDB: func() *mockDB {
				return &mockDB{
					error: gorm.ErrRecordNotFound,
				}
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "Empty username",
			username: "",
			mockDB: func() *mockDB {
				return &mockDB{
					error: gorm.ErrRecordNotFound,
				}
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "Database connection error",
			username: "testuser",
			mockDB: func() *mockDB {
				return &mockDB{
					error: errors.New("database connection error"),
				}
			},
			want:    nil,
			wantErr: errors.New("database connection error"),
		},
		{
			name:     "Special characters in username",
			username: "test@user#123",
			mockDB: func() *mockDB {
				return &mockDB{
					user: &model.User{
						Username: "test@user#123",
						Email:    "special@example.com",
						Bio:      "Special user",
						Image:    "special.jpg",
					},
				}
			},
			want: &model.User{
				Username: "test@user#123",
				Email:    "special@example.com",
				Bio:      "Special user",
				Image:    "special.jpg",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.mockDB()
			store := &UserStore{
				db: mockDB,
			}

			got, err := store.GetByUsername(tt.username)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Bio, got.Bio)
				assert.Equal(t, tt.want.Image, got.Image)
			}
		})
	}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{
		Error: m.error,
		Value: m.user,
	}
}

/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c


 */
func (m *mockDBIsFollowing) Count(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *mockDBIsFollowing) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}

func TestIsFollowing(t *testing.T) {
	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		setupMock   func(*mockDBIsFollowing)
		expected    bool
		expectedErr error
	}{
		{
			name:  "Valid Users - A following B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMock: func(m *mockDBIsFollowing) {
				db := &gorm.DB{Error: nil}
				m.On("Table", "follows").Return(db)
				m.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(2)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					count := args.Get(0).(*int)
					*count = 1
				}).Return(&gorm.DB{Error: nil})
			},
			expected:    true,
			expectedErr: nil,
		},
		{
			name:  "Valid Users - A not following B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMock: func(m *mockDBIsFollowing) {
				db := &gorm.DB{Error: nil}
				m.On("Table", "follows").Return(db)
				m.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(2)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					count := args.Get(0).(*int)
					*count = 0
				}).Return(&gorm.DB{Error: nil})
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User A",
			userA:       nil,
			userB:       &model.User{Model: gorm.Model{ID: 2}},
			setupMock:   func(m *mockDBIsFollowing) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User B",
			userA:       &model.User{Model: gorm.Model{ID: 1}},
			userB:       nil,
			setupMock:   func(m *mockDBIsFollowing) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Both Users Nil",
			userA:       nil,
			userB:       nil,
			setupMock:   func(m *mockDBIsFollowing) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:  "Database Error",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMock: func(m *mockDBIsFollowing) {
				db := &gorm.DB{Error: errors.New("database error")}
				m.On("Table", "follows").Return(db)
				m.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(2)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Return(db)
			},
			expected:    false,
			expectedErr: errors.New("database error"),
		},
		{
			name:  "Same User Reference",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(m *mockDBIsFollowing) {
				db := &gorm.DB{Error: nil}
				m.On("Table", "follows").Return(db)
				m.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(1)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					count := args.Get(0).(*int)
					*count = 0
				}).Return(&gorm.DB{Error: nil})
			},
			expected:    false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mockDBIsFollowing)
			tt.setupMock(mockDB)

			store := &UserStore{
				db: mockDB,
			}

			result, err := store.IsFollowing(tt.userA, tt.userB)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *mockDBIsFollowing) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := append([]interface{}{query}, args...)
	return m.Called(callArgs...).Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9


 */
func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}{
		{
			name:     "Successful UserStore creation with valid DB",
			db:       &gorm.DB{},
			wantNil:  false,
			scenario: "Scenario 1: Basic initialization with valid DB",
		},
		{
			name:     "UserStore creation with nil DB",
			db:       nil,
			wantNil:  false,
			scenario: "Scenario 2: Handling nil DB parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Testing scenario:", tt.scenario)

			userStore := NewUserStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, userStore, "UserStore should be nil")
			} else {
				assert.NotNil(t, userStore, "UserStore should not be nil")
				assert.Equal(t, tt.db, userStore.db, "DB reference should match input")
			}
		})
	}

	t.Run("Verify DB reference integrity", func(t *testing.T) {
		mockDB := &gorm.DB{}
		userStore := NewUserStore(mockDB)
		assert.Same(t, mockDB, userStore.db, "DB reference should be the same instance")
	})

	t.Run("Create multiple independent instances", func(t *testing.T) {
		db1 := &gorm.DB{}
		db2 := &gorm.DB{}

		store1 := NewUserStore(db1)
		store2 := NewUserStore(db2)

		assert.NotEqual(t, store1, store2, "Different instances should not be equal")
		assert.NotSame(t, store1.db, store2.db, "DB references should be independent")
	})

	t.Run("Concurrent UserStore creation", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		stores := make([]*UserStore, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				mockDB := &gorm.DB{}
				stores[index] = NewUserStore(mockDB)
			}(i)
		}

		wg.Wait()

		for i, store := range stores {
			assert.NotNil(t, store, "Store %d should not be nil", i)
			assert.NotNil(t, store.db, "DB in store %d should not be nil", i)
		}
	})
}

/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55


 */
func TestUnfollow(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*testing.T) (*UserStore, *model.User, *model.User)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful unfollow",
			setup: func(t *testing.T) (*UserStore, *model.User, *model.User) {
				db, mock, err := setupTestDB(t)
				if err != nil {
					t.Fatalf("failed to setup test db: %v", err)
				}

				userA := &model.User{
					Model:    gorm.Model{ID: 1},
					Username: "userA",
					Email:    "userA@test.com",
				}
				userB := &model.User{
					Model:    gorm.Model{ID: 2},
					Username: "userB",
					Email:    "userB@test.com",
				}

				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(userA.ID, userB.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				return &UserStore{db: db}, userA, userB
			},
			wantErr: false,
		},
		{
			name: "Unfollow non-followed user",
			setup: func(t *testing.T) (*UserStore, *model.User, *model.User) {
				db, mock, err := setupTestDB(t)
				if err != nil {
					t.Fatalf("failed to setup test db: %v", err)
				}

				userA := &model.User{
					Model:    gorm.Model{ID: 1},
					Username: "userA",
					Email:    "userA@test.com",
				}
				userB := &model.User{
					Model:    gorm.Model{ID: 2},
					Username: "userB",
					Email:    "userB@test.com",
				}

				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(userA.ID, userB.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()

				return &UserStore{db: db}, userA, userB
			},
			wantErr: false,
		},
		{
			name: "Invalid user reference",
			setup: func(t *testing.T) (*UserStore, *model.User, *model.User) {
				db, _, err := setupTestDB(t)
				if err != nil {
					t.Fatalf("failed to setup test db: %v", err)
				}

				userA := &model.User{
					Model:    gorm.Model{ID: 1},
					Username: "userA",
					Email:    "userA@test.com",
				}

				return &UserStore{db: db}, userA, nil
			},
			wantErr: true,
			errMsg:  "invalid user reference",
		},
		{
			name: "Database connection error",
			setup: func(t *testing.T) (*UserStore, *model.User, *model.User) {
				db, mock, err := setupTestDB(t)
				if err != nil {
					t.Fatalf("failed to setup test db: %v", err)
				}

				userA := &model.User{
					Model:    gorm.Model{ID: 1},
					Username: "userA",
					Email:    "userA@test.com",
				}
				userB := &model.User{
					Model:    gorm.Model{ID: 2},
					Username: "userB",
					Email:    "userB@test.com",
				}

				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `follows`").
					WithArgs(userA.ID, userB.ID).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()

				return &UserStore{db: db}, userA, userB
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, userA, userB := tt.setup(t)
			err := store.Unfollow(userA, userB)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnfollowConcurrent(t *testing.T) {
	store, userA, userB := setupConcurrentTest(t)

	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			err := store.Unfollow(userA, userB)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
}

func setupConcurrentTest(t *testing.T) (*UserStore, *model.User, *model.User) {
	db, mock, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	userA := &model.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "userA",
		Email:    "userA@test.com",
	}

	userB := &model.User{
		Model: gorm.Model{
			ID:        2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: "userB",
		Email:    "userB@test.com",
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM follows").
		WithArgs(userA.ID, userB.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	return &UserStore{db: db}, userA, userB
}

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435


 */
func TestUpdate(t *testing.T) {
	tests := []struct {
		name      string
		user      *model.User
		setupMock func(sqlmock.Sqlmock)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "Successful Update",
			user: &model.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Username: "updated_user",
				Email:    "updated@example.com",
				Password: "newpassword",
				Bio:      "Updated bio",
				Image:    "new_image.jpg",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WithArgs(
						sqlmock.AnyArg(),
						"updated_user",
						"updated@example.com",
						"newpassword",
						"Updated bio",
						"new_image.jpg",
						1,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Non-Existent User",
			user: &model.User{
				Model: gorm.Model{ID: 999},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "record not found",
		},
		{
			name: "Database Connection Error",
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(errors.New("database connection error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database connection error",
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

			tt.setupMock(mock)

			store := &UserStore{db: gormDB}
			err = store.Update(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7


 */
func TestGetFollowingUserIDs(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		mockSetup func(sqlmock.Sqlmock)
		expected  []uint
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "Successfully retrieve following user IDs",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(uint(2)).
					AddRow(uint(3)).
					AddRow(uint(4))
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			expected: []uint{2, 3, 4},
			wantErr:  false,
		},
		{
			name:   "User with no followings",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			expected: []uint{},
			wantErr:  false,
		},
		{
			name:   "Database connection error",
			userID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(uint(1)).
					WillReturnError(errors.New("database connection error"))
			},
			expected: []uint{},
			wantErr:  true,
			errMsg:   "database connection error",
		},
		{
			name:   "Invalid user ID",
			userID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("^SELECT to_user_id FROM follows WHERE").
					WithArgs(uint(999)).
					WillReturnRows(rows)
			},
			expected: []uint{},
			wantErr:  false,
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

			store := &UserStore{
				db: gormDB,
			}

			user := &model.User{
				Model: gorm.Model{ID: tt.userID},
			}

			got, err := store.GetFollowingUserIDs(user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, got)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

