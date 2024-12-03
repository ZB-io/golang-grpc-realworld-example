package store

import (
	"errors"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sync"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"database/sql"
)

/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920


 */
func TestUserCreate(t *testing.T) {
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
				Image:    "https://example.com/image.jpg",
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Duplicate username",
			user: &model.User{
				Username: "existinguser",
				Email:    "another@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			dbError: errors.New("duplicate key value violates unique constraint"),
			wantErr: true,
		},
		{
			name: "Missing required fields",
			user: &model.User{
				Username: "",
				Email:    "",
				Password: "",
			},
			dbError: errors.New("not null constraint violation"),
			wantErr: true,
		},
		{
			name: "Database connection failure",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			dbError: errors.New("connection refused"),
			wantErr: true,
		},
		{
			name: "Large data fields",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Bio:      string(make([]byte, 1000)),
				Image:    string(make([]byte, 1000)),
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Special characters in fields",
			user: &model.User{
				Username: "test@user#$%",
				Email:    "test+special@example.com",
				Password: "password!@#$%^&*()",
				Bio:      "Bio with émojis 🎉",
				Image:    "https://example.com/image?special=true&param=value",
			},
			dbError: nil,
			wantErr: false,
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

			if tt.dbError != nil {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnError(tt.dbError)
				mock.ExpectRollback()
			} else {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

			err = store.Create(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				t.Logf("Expected error occurred: %v", err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.user.ID)
				assert.False(t, tt.user.CreatedAt.IsZero())
				t.Log("User created successfully")
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
func NewMockDB() (*MockDB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		return nil, err
	}

	return &MockDB{
		DB:   gormDB,
		mock: mock,
	}, nil
}

func TestFollow(t *testing.T) {
	tests := []struct {
		name    string
		userA   *model.User
		userB   *model.User
		dbSetup func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful follow operation",
			userA: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "userA",
				Email:    "userA@test.com",
				Password: "password",
				Bio:      "test bio",
				Image:    "test.jpg",
			},
			userB: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "userB",
				Email:    "userB@test.com",
				Password: "password",
				Bio:      "test bio",
				Image:    "test.jpg",
			},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:    "Nil user A",
			userA:   nil,
			userB:   &model.User{Model: gorm.Model{ID: 2}},
			wantErr: true,
			errMsg:  "invalid user parameters",
		},
		{
			name:    "Nil user B",
			userA:   &model.User{Model: gorm.Model{ID: 1}},
			userB:   nil,
			wantErr: true,
			errMsg:  "invalid user parameters",
		},
		{
			name:    "Self follow attempt",
			userA:   &model.User{Model: gorm.Model{ID: 1}},
			userB:   &model.User{Model: gorm.Model{ID: 1}},
			wantErr: true,
			errMsg:  "self-follow not allowed",
		},
		{
			name:  "Already following user",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("duplicate entry"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "duplicate entry",
		},
		{
			name:  "Database error",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			dbSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `follows`").
					WithArgs(1, 2).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := NewMockDB()
			if err != nil {
				t.Fatalf("failed to create mock db: %v", err)
			}
			defer mockDB.DB.Close()

			if tt.dbSetup != nil {
				tt.dbSetup(mockDB.mock)
			}

			store := &UserStore{
				db: mockDB.DB,
			}

			err = store.Follow(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("Follow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Follow() error message = %v, want %v", err.Error(), tt.errMsg)
			}

			if err := mockDB.mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1


 */
func (m *UserMockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	called := m.Called(out, where)
	return called.Get(0).(*gorm.DB)
}

func TestGetByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		mockSetup     func(*UserMockDB)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:  "Successful user retrieval",
			email: "test@example.com",
			mockSetup: func(mock *UserMockDB) {
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
					Return(&gorm.DB{Error: nil}).
					Run(func(args mock.Arguments) {
						arg := args.Get(0).(*model.User)
						*arg = *expectedUser
					})
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Email:    "test@example.com",
				Username: "testuser",
			},
			expectedError: nil,
		},
		{
			name:  "User not found",
			email: "nonexistent@example.com",
			mockSetup: func(mock *UserMockDB) {
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
			mockSetup: func(mock *UserMockDB) {
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
			mockSetup: func(mock *UserMockDB) {
				dbError := errors.New("database connection error")
				mock.On("Where", "email = ?", []interface{}{"test@example.com"}).
					Return(&gorm.DB{Error: dbError})
				mock.On("First", mock.Anything, mock.Anything).
					Return(&gorm.DB{Error: dbError})
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(UserMockDB)
			tt.mockSetup(mockDB)

			store := &UserStore{
				db: &gorm.DB{},
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

func (m *UserMockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	called := m.Called(query, args)
	return called.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1


 */
func (m *UserMockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		setupMock   func(*UserMockDB)
		expectUser  *model.User
		expectError error
	}{
		{
			name: "Successfully retrieve user by valid ID",
			id:   1,
			setupMock: func(m *UserMockDB) {
				expectedUser := &model.User{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Username: "testuser",
					Email:    "test@example.com",
					Bio:      "Test bio",
					Image:    "test-image.jpg",
				}
				m.On("Find", mock.AnythingOfType("*model.User"), uint(1)).Return(
					&gorm.DB{Error: nil},
				).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.User)
					*arg = *expectedUser
				})
			},
			expectUser: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "testuser",
				Email:    "test@example.com",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			expectError: nil,
		},
		{
			name: "Attempt to retrieve non-existent user ID",
			id:   999,
			setupMock: func(m *UserMockDB) {
				m.On("Find", mock.AnythingOfType("*model.User"), uint(999)).Return(
					&gorm.DB{Error: gorm.ErrRecordNotFound},
				)
			},
			expectUser:  nil,
			expectError: gorm.ErrRecordNotFound,
		},
		{
			name: "Handle database connection error",
			id:   1,
			setupMock: func(m *UserMockDB) {
				m.On("Find", mock.AnythingOfType("*model.User"), uint(1)).Return(
					&gorm.DB{Error: errors.New("database connection error")},
				)
			},
			expectUser:  nil,
			expectError: errors.New("database connection error"),
		},
		{
			name: "Handle zero ID value",
			id:   0,
			setupMock: func(m *UserMockDB) {
				m.On("Find", mock.AnythingOfType("*model.User"), uint(0)).Return(
					&gorm.DB{Error: gorm.ErrRecordNotFound},
				)
			},
			expectUser:  nil,
			expectError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(UserMockDB)
			tt.setupMock(mockDB)

			store := &UserStore{
				db: &gorm.DB{},
			}

			user, err := store.GetByID(tt.id)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectUser.ID, user.ID)
				assert.Equal(t, tt.expectUser.Username, user.Username)
				assert.Equal(t, tt.expectUser.Email, user.Email)
				assert.Equal(t, tt.expectUser.Bio, user.Bio)
				assert.Equal(t, tt.expectUser.Image, user.Image)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24


 */
func (m *UserMockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	called := m.Called(out, where)
	return called.Get(0).(*gorm.DB)
}

func TestGetByUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockSetup     func(*UserMockDB)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:     "Successful user retrieval",
			username: "validuser",
			mockSetup: func(mock *UserMockDB) {
				expectedUser := &model.User{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Username: "validuser",
					Email:    "valid@example.com",
					Bio:      "Test bio",
				}

				db := &gorm.DB{Error: nil}
				mock.On("Where", "username = ?", []interface{}{"validuser"}).Return(db)
				mock.On("First", mock.Anything, mock.Anything).Return(db).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.User)
					*arg = *expectedUser
				})
			},
			expectedUser: &model.User{
				Username: "validuser",
				Email:    "valid@example.com",
				Bio:      "Test bio",
			},
			expectedError: nil,
		},
		{
			name:     "Non-existent username",
			username: "nonexistent",
			mockSetup: func(mock *UserMockDB) {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				mock.On("Where", "username = ?", []interface{}{"nonexistent"}).Return(db)
				mock.On("First", mock.Anything, mock.Anything).Return(db)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Empty username",
			username: "",
			mockSetup: func(mock *UserMockDB) {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				mock.On("Where", "username = ?", []interface{}{""}).Return(db)
				mock.On("First", mock.Anything, mock.Anything).Return(db)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Database connection error",
			username: "validuser",
			mockSetup: func(mock *UserMockDB) {
				db := &gorm.DB{Error: errors.New("database connection failed")}
				mock.On("Where", "username = ?", []interface{}{"validuser"}).Return(db)
				mock.On("First", mock.Anything, mock.Anything).Return(db)
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(UserMockDB)
			tt.mockSetup(mockDB)

			store := &UserStore{
				db: mockDB,
			}

			user, err := store.GetByUsername(tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *UserMockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	called := m.Called(query, args)
	return called.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c


 */
func (m *UserMockDBForFollowing) Count(value interface{}) *gorm.DB {
	return m.Called(value).Get(0).(*gorm.DB)
}

func (m *UserMockDBForFollowing) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}

func TestIsFollowing(t *testing.T) {
	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		setupMock   func(*UserMockDBForFollowing)
		wantResult  bool
		wantErr     error
		description string
	}{
		{
			name:  "Valid Users - User A Following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMock: func(m *UserMockDBForFollowing) {
				db := &gorm.DB{}
				m.On("Table", "follows").Return(db)
				m.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(2)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					count := args.Get(0).(*int)
					*count = 1
				}).Return(&gorm.DB{Error: nil})
			},
			wantResult:  true,
			wantErr:     nil,
			description: "Tests basic functionality where User A is following User B",
		},
		{
			name:  "Valid Users - User A Not Following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMock: func(m *UserMockDBForFollowing) {
				db := &gorm.DB{}
				m.On("Table", "follows").Return(db)
				m.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(2)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					count := args.Get(0).(*int)
					*count = 0
				}).Return(&gorm.DB{Error: nil})
			},
			wantResult:  false,
			wantErr:     nil,
			description: "Tests case where User A is not following User B",
		},
		{
			name:        "Nil User A Parameter",
			userA:       nil,
			userB:       &model.User{Model: gorm.Model{ID: 2}},
			setupMock:   func(m *UserMockDBForFollowing) {},
			wantResult:  false,
			wantErr:     nil,
			description: "Tests behavior when first user parameter is nil",
		},
		{
			name:        "Nil User B Parameter",
			userA:       &model.User{Model: gorm.Model{ID: 1}},
			userB:       nil,
			setupMock:   func(m *UserMockDBForFollowing) {},
			wantResult:  false,
			wantErr:     nil,
			description: "Tests behavior when second user parameter is nil",
		},
		{
			name:  "Database Error",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMock: func(m *UserMockDBForFollowing) {
				db := &gorm.DB{Error: errors.New("database error")}
				m.On("Table", "follows").Return(db)
				m.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(2)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Return(db)
			},
			wantResult:  false,
			wantErr:     errors.New("database error"),
			description: "Tests behavior when database query returns an error",
		},
		{
			name:        "Both Users Are Nil",
			userA:       nil,
			userB:       nil,
			setupMock:   func(m *UserMockDBForFollowing) {},
			wantResult:  false,
			wantErr:     nil,
			description: "Tests edge case where both user parameters are nil",
		},
		{
			name:  "Same User Reference",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 1}},
			setupMock: func(m *UserMockDBForFollowing) {
				db := &gorm.DB{}
				m.On("Table", "follows").Return(db)
				m.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(1)).Return(db)
				m.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					count := args.Get(0).(*int)
					*count = 0
				}).Return(&gorm.DB{Error: nil})
			},
			wantResult:  false,
			wantErr:     nil,
			description: "Tests case where same user is passed as both parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)

			mockDB := new(UserMockDBForFollowing)
			tt.setupMock(mockDB)

			store := &UserStore{
				db: mockDB,
			}

			got, err := store.IsFollowing(tt.userA, tt.userB)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantResult, got)

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *UserMockDBForFollowing) Where(query interface{}, args ...interface{}) *gorm.DB {
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
			name: "Scenario 1: Successfully Create New UserStore with Valid DB Connection",
			db: &gorm.DB{
				Error: nil,
			},
			wantNil:  false,
			scenario: "Basic initialization with valid DB",
		},
		{
			name:     "Scenario 2: Create UserStore with Nil DB Parameter",
			db:       nil,
			wantNil:  false,
			scenario: "Initialization with nil DB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting test:", tt.scenario)

			userStore := NewUserStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, userStore, "UserStore should be nil")
			} else {
				assert.NotNil(t, userStore, "UserStore should not be nil")
				assert.Equal(t, tt.db, userStore.db, "DB reference should match input")
			}

			t.Log("Test completed successfully")
		})
	}

	t.Run("Scenario 3: Verify DB Reference Integrity", func(t *testing.T) {
		mockDB := &gorm.DB{
			Error: nil,
		}
		userStore := NewUserStore(mockDB)
		assert.Same(t, mockDB, userStore.db, "DB reference should be the same instance")
	})

	t.Run("Scenario 4: Create Multiple UserStore Instances", func(t *testing.T) {
		db1 := &gorm.DB{}
		db2 := &gorm.DB{}
		store1 := NewUserStore(db1)
		store2 := NewUserStore(db2)
		assert.NotEqual(t, store1, store2, "Different instances should not be equal")
		assert.Equal(t, db1, store1.db, "First store should have correct DB reference")
		assert.Equal(t, db2, store2.db, "Second store should have correct DB reference")
	})

	t.Run("Scenario 7: Concurrent UserStore Creation", func(t *testing.T) {
		const numGoroutines = 10
		var wg sync.WaitGroup
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

		for i := 0; i < numGoroutines; i++ {
			assert.NotNil(t, stores[i], "Concurrent creation should succeed")
			assert.NotNil(t, stores[i].db, "DB reference should not be nil")
		}
	})
}

/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55


 */
func TestUnfollow(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(*gorm.DB) (*model.User, *model.User)
		mockDBError error
		wantErr     bool
		errMsg      string
	}{
		{
			name: "Successful unfollow",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userA", Email: "userA@test.com"}
				userB := &model.User{Username: "userB", Email: "userB@test.com"}
				db.Create(userA)
				db.Create(userB)
				db.Model(userA).Association("Follows").Append(userB)
				return userA, userB
			},
			wantErr: false,
		},
		{
			name: "Unfollow non-followed user",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userC", Email: "userC@test.com"}
				userB := &model.User{Username: "userD", Email: "userD@test.com"}
				db.Create(userA)
				db.Create(userB)
				return userA, userB
			},
			wantErr: false,
		},
		{
			name: "Nil user parameters",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userE", Email: "userE@test.com"}
				db.Create(userA)
				return userA, nil
			},
			wantErr: true,
			errMsg:  "invalid user parameters",
		},
		{
			name: "Database error",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userF", Email: "userF@test.com"}
				userB := &model.User{Username: "userG", Email: "userG@test.com"}
				return userA, userB
			},
			mockDBError: errors.New("database connection error"),
			wantErr:     true,
			errMsg:      "database connection error",
		},
		{
			name: "Soft deleted user",
			setupFunc: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userH", Email: "userH@test.com"}
				userB := &model.User{Username: "userI", Email: "userI@test.com"}
				db.Create(userA)
				db.Create(userB)
				db.Model(userA).Association("Follows").Append(userB)
				db.Delete(userB)
				return userA, userB
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, err := setupUnfollowTestDB()
			assert.NoError(t, err)
			defer mockDB.Close()

			if tt.mockDBError != nil {
				mockDB.Close()
			}

			store := &UserStore{db: mockDB}
			userA, userB := tt.setupFunc(mockDB)

			if tt.name == "Successful unfollow" {
				var wg sync.WaitGroup
				for i := 0; i < 5; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						err := store.Unfollow(userA, userB)
						assert.NoError(t, err)
					}()
				}
				wg.Wait()
			}

			err = store.Unfollow(userA, userB)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				var follows []model.User
				mockDB.Model(userA).Association("Follows").Find(&follows)
				for _, follow := range follows {
					assert.NotEqual(t, userB.ID, follow.ID)
				}
			}

			t.Logf("Test case '%s' completed", tt.name)
		})
	}
}

func setupUnfollowTestDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.User{})

	return db, nil
}

/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435


 */
func (m *UserMockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		setupFn func(*UserMockDB)
		wantErr error
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
				Image:    "updated.jpg",
			},
			setupFn: func(mockDB *UserMockDB) {
				mockDB.On("Model", mock.Anything).Return(&gorm.DB{Error: nil})
				mockDB.On("Update", mock.Anything).Return(&gorm.DB{Error: nil})
			},
			wantErr: nil,
		},
		{
			name: "Update with Duplicate Username",
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "existing_user",
			},
			setupFn: func(mockDB *UserMockDB) {
				mockDB.On("Model", mock.Anything).Return(&gorm.DB{Error: nil})
				mockDB.On("Update", mock.Anything).Return(&gorm.DB{
					Error: errors.New("duplicate key value violates unique constraint"),
				})
			},
			wantErr: errors.New("duplicate key value violates unique constraint"),
		},
		{
			name: "Update Non-Existent User",
			user: &model.User{
				Model: gorm.Model{ID: 999},
			},
			setupFn: func(mockDB *UserMockDB) {
				mockDB.On("Model", mock.Anything).Return(&gorm.DB{Error: nil})
				mockDB.On("Update", mock.Anything).Return(&gorm.DB{
					Error: gorm.ErrRecordNotFound,
				})
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			setupFn: func(mockDB *UserMockDB) {
				mockDB.On("Model", mock.Anything).Return(&gorm.DB{Error: nil})
				mockDB.On("Update", mock.Anything).Return(&gorm.DB{
					Error: errors.New("database connection error"),
				})
			},
			wantErr: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(UserMockDB)
			tt.setupFn(mockDB)

			store := &UserStore{
				db: &gorm.DB{},
			}

			err := store.Update(tt.user)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *UserMockDB) Update(attrs interface{}) *gorm.DB {
	args := m.Called(attrs)
	return args.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7


 */
func TestGetFollowingUserIDs_Store(t *testing.T) {
	tests := []struct {
		name      string
		user      *model.User
		mockSetup func(sqlmock.Sqlmock)
		expected  []uint
		wantErr   bool
		errMsg    string
	}{
		{
			name: "Successfully retrieve following user IDs",
			user: &model.User{Model: gorm.Model{ID: 1}},
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
			name: "User with no followings",
			user: &model.User{Model: gorm.Model{ID: 1}},
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
			name: "Database connection error",
			user: &model.User{Model: gorm.Model{ID: 1}},
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
			name: "Invalid user ID",
			user: &model.User{Model: gorm.Model{ID: 999}},
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

			store := &UserStore{db: gormDB}
			got, err := store.GetFollowingUserIDs(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
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

