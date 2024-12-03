package store

import (
	"errors"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
	"database/sql"
	"sync"
)

/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920


 */
func TestCreate(t *testing.T) {

	tests := []struct {
		name    string
		user    *model.User
		setupDB func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
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
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Duplicate username",
			user: &model.User{
				Username: "existinguser",
				Email:    "new@example.com",
				Password: "password123",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(errors.New("Error 1062: Duplicate entry"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "Duplicate entry",
		},
		{
			name: "Missing required fields",
			user: &model.User{
				Username: "",
				Email:    "",
			},
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(errors.New("not null constraint failed"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "not null constraint failed",
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
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnError(errors.New("connection refused"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "connection refused",
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
			setupDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Special characters in fields",
			user: &model.User{
				Username: "test@user#$%",
				Email:    "test+special@example.com",
				Password: "pass!@#$%^&*()",
				Bio:      "Bio with émojis 🎉",
				Image:    "https://example.com/image?special=true&param=value",
			},
			setupDB: func(mock sqlmock.Sqlmock) {
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

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open gorm DB: %v", err)
			}
			defer gormDB.Close()

			tt.setupDB(mock)

			store := &UserStore{db: gormDB}

			err = store.Create(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				t.Logf("Expected error received: %v", err)
			} else {
				assert.NoError(t, err)
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
func TestFollow(t *testing.T) {

	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	store := &UserStore{db: db}

	tests := []struct {
		name    string
		setup   func(t *testing.T) (*model.User, *model.User)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful follow",
			setup: func(t *testing.T) (*model.User, *model.User) {
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
				require.NoError(t, db.Create(userA).Error)
				require.NoError(t, db.Create(userB).Error)
				return userA, userB
			},
			wantErr: false,
		},
		{
			name: "Nil user parameters",
			setup: func(t *testing.T) (*model.User, *model.User) {
				return nil, nil
			},
			wantErr: true,
			errMsg:  "invalid user parameters",
		},
		{
			name: "Self follow attempt",
			setup: func(t *testing.T) (*model.User, *model.User) {
				user := &model.User{
					Username: "selfFollow",
					Email:    "self@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				require.NoError(t, db.Create(user).Error)
				return user, user
			},
			wantErr: false,
		},
		{
			name: "Already following user",
			setup: func(t *testing.T) (*model.User, *model.User) {
				userA := &model.User{
					Username: "duplicateA",
					Email:    "dupA@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				userB := &model.User{
					Username: "duplicateB",
					Email:    "dupB@test.com",
					Password: "password",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				require.NoError(t, db.Create(userA).Error)
				require.NoError(t, db.Create(userB).Error)
				require.NoError(t, store.Follow(userA, userB))
				return userA, userB
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cleanupDB(t, db)

			userA, userB := tt.setup(t)
			err := store.Follow(userA, userB)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

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
		})
	}
}

func cleanupDB(t *testing.T, db *gorm.DB) {

	t.Helper()
	db.Exec("DELETE FROM follows")
	db.Exec("DELETE FROM users")
}

func setupTestDB() (*gorm.DB, error) {

	return nil, nil
}

/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1


 */
func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	called := m.Called(out, where)
	return called.Get(0).(*gorm.DB)
}

func TestGetByEmail(t *testing.T) {

	tests := []struct {
		name          string
		email         string
		mockSetup     func(*MockDB)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:  "Successfully retrieve user by valid email",
			email: "test@example.com",
			mockSetup: func(mock *MockDB) {
				expectedUser := &model.User{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Email:    "test@example.com",
					Username: "testuser",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				db := &gorm.DB{Error: nil}
				mock.On("Where", "email = ?", []interface{}{"test@example.com"}).Return(db)
				mock.On("First", mock.Anything, mock.Anything).Return(db).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.User)
					*arg = *expectedUser
				})
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Email:    "test@example.com",
				Username: "testuser",
				Bio:      "test bio",
				Image:    "test.jpg",
			},
			expectedError: nil,
		},
		{
			name:  "Non-existent email",
			email: "nonexistent@example.com",
			mockSetup: func(mock *MockDB) {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				mock.On("Where", "email = ?", []interface{}{"nonexistent@example.com"}).Return(db)
				mock.On("First", mock.Anything, mock.Anything).Return(db)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Empty email",
			email: "",
			mockSetup: func(mock *MockDB) {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				mock.On("Where", "email = ?", []interface{}{""}).Return(db)
				mock.On("First", mock.Anything, mock.Anything).Return(db)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Database connection error",
			email: "test@example.com",
			mockSetup: func(mock *MockDB) {
				db := &gorm.DB{Error: errors.New("database connection error")}
				mock.On("Where", "email = ?", []interface{}{"test@example.com"}).Return(db)
				mock.On("First", mock.Anything, mock.Anything).Return(db)
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
		{
			name:  "Special characters in email",
			email: "test+label@example.com",
			mockSetup: func(mock *MockDB) {
				expectedUser := &model.User{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Email:    "test+label@example.com",
					Username: "testuser",
					Bio:      "test bio",
					Image:    "test.jpg",
				}
				db := &gorm.DB{Error: nil}
				mock.On("Where", "email = ?", []interface{}{"test+label@example.com"}).Return(db)
				mock.On("First", mock.Anything, mock.Anything).Return(db).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.User)
					*arg = *expectedUser
				})
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Email:    "test+label@example.com",
				Username: "testuser",
				Bio:      "test bio",
				Image:    "test.jpg",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(MockDB)
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
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
				assert.Equal(t, tt.expectedUser.Image, user.Image)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	called := m.Called(query, args)
	return called.Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1


 */
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func TestGetByID(t *testing.T) {

	tests := []struct {
		name          string
		userID        uint
		setupMock     func(*MockDB) *gorm.DB
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:   "Successfully retrieve user by valid ID",
			userID: 1,
			setupMock: func(m *MockDB) *gorm.DB {
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
				db := &gorm.DB{Error: nil}
				m.On("Find", mock.AnythingOfType("*model.User"), uint(1)).Return(db)

				if u, ok := mock.AnythingOfType("*model.User").(*model.User); ok {
					*u = *expectedUser
				}
				return db
			},
			expectedUser: &model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Username: "testuser",
				Email:    "test@example.com",
				Bio:      "Test bio",
				Image:    "test-image.jpg",
			},
			expectedError: nil,
		},
		{
			name:   "Attempt to retrieve non-existent user ID",
			userID: 999,
			setupMock: func(m *MockDB) *gorm.DB {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				m.On("Find", mock.AnythingOfType("*model.User"), uint(999)).Return(db)
				return db
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:   "Handle database connection error",
			userID: 1,
			setupMock: func(m *MockDB) *gorm.DB {
				db := &gorm.DB{Error: errors.New("database connection failed")}
				m.On("Find", mock.AnythingOfType("*model.User"), uint(1)).Return(db)
				return db
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection failed"),
		},
		{
			name:   "Handle zero ID value",
			userID: 0,
			setupMock: func(m *MockDB) *gorm.DB {
				db := &gorm.DB{Error: gorm.ErrRecordNotFound}
				m.On("Find", mock.AnythingOfType("*model.User"), uint(0)).Return(db)
				return db
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(MockDB)
			db := tt.setupMock(mockDB)

			store := &UserStore{
				db: db,
			}

			user, err := store.GetByID(tt.userID)

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
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
				assert.Equal(t, tt.expectedUser.Image, user.Image)
			}

			mockDB.AssertExpectations(t)
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
		username      string
		mockSetup     func(sqlmock.Sqlmock)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:     "Successfully retrieve existing user",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"username", "email", "password", "bio", "image",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"testuser", "test@example.com", "hashedpass", "test bio", "image.jpg",
				)
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE").
					WithArgs("testuser").
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
			name:     "User not found",
			username: "nonexistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE").
					WithArgs("nonexistent").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Empty username",
			username: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE").
					WithArgs("").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Database connection error",
			username: "testuser",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE").
					WithArgs("testuser").
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
		{
			name:     "Special characters in username",
			username: "test@user#123",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "created_at", "updated_at", "deleted_at",
					"username", "email", "password", "bio", "image",
				}).AddRow(
					1, time.Now(), time.Now(), nil,
					"test@user#123", "special@example.com", "hashedpass", "test bio", "image.jpg",
				)
				mock.ExpectQuery("^SELECT (.+) FROM \"users\" WHERE").
					WithArgs("test@user#123").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "test@user#123",
				Email:    "special@example.com",
				Password: "hashedpass",
				Bio:      "test bio",
				Image:    "image.jpg",
			},
			expectedError: nil,
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

			user, err := store.GetByUsername(tt.username)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Bio, user.Bio)
				assert.Equal(t, tt.expectedUser.Image, user.Image)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c


 */
func (m *MockDB) Count(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}

func TestIsFollowing(t *testing.T) {

	tests := []struct {
		name        string
		userA       *model.User
		userB       *model.User
		setupMock   func(*MockDB)
		expected    bool
		expectedErr error
	}{
		{
			name:  "Valid Users - A following B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMock: func(m *MockDB) {
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
			setupMock: func(m *MockDB) {
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
			setupMock:   func(m *MockDB) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Nil User B",
			userA:       &model.User{Model: gorm.Model{ID: 1}},
			userB:       nil,
			setupMock:   func(m *MockDB) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:        "Both Users Nil",
			userA:       nil,
			userB:       nil,
			setupMock:   func(m *MockDB) {},
			expected:    false,
			expectedErr: nil,
		},
		{
			name:  "Database Error",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			setupMock: func(m *MockDB) {
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
			setupMock: func(m *MockDB) {
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

			mockDB := new(MockDB)
			tt.setupMock(mockDB)

			store := &UserStore{
				db: mockDB,
			}

			t.Logf("Running test case: %s", tt.name)
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

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := append([]interface{}{query}, args...)
	return m.Called(callArgs...).Get(0).(*gorm.DB)
}

/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9


 */
func TestNewUserStore(t *testing.T) {

	createMockDB := func(identifier string) *gorm.DB {
		db := &gorm.DB{
			Value:     identifier,
			LogMode:   true,
			callbacks: &gorm.Callback{},
		}
		return db
	}

	tests := []struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		scenario string
	}{
		{
			name:     "Scenario 1: Successfully Create New UserStore with Valid DB Connection",
			db:       createMockDB("valid-db"),
			wantNil:  false,
			scenario: "Basic initialization with valid DB",
		},
		{
			name:     "Scenario 2: Create UserStore with Nil DB Parameter",
			db:       nil,
			wantNil:  false,
			scenario: "Handling nil DB parameter",
		},
		{
			name:     "Scenario 3: Verify DB Reference Integrity",
			db:       createMockDB("reference-check"),
			wantNil:  false,
			scenario: "DB reference maintenance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Starting:", tt.scenario)

			userStore := NewUserStore(tt.db)

			if tt.wantNil {
				assert.Nil(t, userStore, "UserStore should be nil")
			} else {
				assert.NotNil(t, userStore, "UserStore should not be nil")
				assert.Equal(t, tt.db, userStore.db, "DB reference should match input")
			}

			t.Log("Completed:", tt.scenario)
		})
	}

	t.Run("Scenario 4: Create Multiple UserStore Instances", func(t *testing.T) {
		db1 := createMockDB("db1")
		db2 := createMockDB("db2")

		store1 := NewUserStore(db1)
		store2 := NewUserStore(db2)

		assert.NotEqual(t, store1, store2, "Different instances should not be equal")
		assert.Equal(t, db1, store1.db, "First store should maintain its DB reference")
		assert.Equal(t, db2, store2.db, "Second store should maintain its DB reference")
	})

	t.Run("Scenario 5: Verify UserStore with Configured DB Settings", func(t *testing.T) {
		db := createMockDB("configured-db")
		db.LogMode = true

		userStore := NewUserStore(db)
		assert.Equal(t, db.LogMode, userStore.db.LogMode, "DB settings should be preserved")
	})

	t.Run("Scenario 7: Concurrent UserStore Creation", func(t *testing.T) {
		const numGoroutines = 10
		var wg sync.WaitGroup
		stores := make([]*UserStore, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				db := createMockDB(fmt.Sprintf("concurrent-db-%d", index))
				stores[index] = NewUserStore(db)
			}(i)
		}

		wg.Wait()

		for i := 0; i < numGoroutines; i++ {
			assert.NotNil(t, stores[i], "Concurrent creation should succeed")
		}
	})
}

/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55


 */
func TestUnfollow(t *testing.T) {

	type testCase struct {
		name    string
		userA   *model.User
		userB   *model.User
		setup   func(*gorm.DB)
		wantErr bool
	}

	createTestUser := func(username string) *model.User {
		return &model.User{
			Model: gorm.Model{
				ID:        uint(time.Now().UnixNano()),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Username: username,
			Email:    username + "@test.com",
			Password: "password",
			Bio:      "test bio",
			Image:    "test-image.jpg",
		}
	}

	tests := []testCase{
		{
			name:  "Successful Unfollow",
			userA: createTestUser("userA"),
			userB: createTestUser("userB"),
			setup: func(db *gorm.DB) {

				db.Model(createTestUser("userA")).Association("Follows").Append(createTestUser("userB"))
			},
			wantErr: false,
		},
		{
			name:  "Unfollow Non-Followed User",
			userA: createTestUser("userC"),
			userB: createTestUser("userD"),
			setup: func(db *gorm.DB) {

			},
			wantErr: false,
		},
		{
			name:    "Nil User A",
			userA:   nil,
			userB:   createTestUser("userE"),
			setup:   func(db *gorm.DB) {},
			wantErr: true,
		},
		{
			name:    "Nil User B",
			userA:   createTestUser("userF"),
			userB:   nil,
			setup:   func(db *gorm.DB) {},
			wantErr: true,
		},
		{
			name:  "Database Error",
			userA: createTestUser("userG"),
			userB: createTestUser("userH"),
			setup: func(db *gorm.DB) {

				db.Close()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			mockDB, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("failed to open mock database: %v", err)
			}
			defer mockDB.Close()

			store := &UserStore{db: mockDB}

			tc.setup(mockDB)

			err = store.Unfollow(tc.userA, tc.userB)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tc.userA != nil && tc.userB != nil {
					var follows []model.User
					err = mockDB.Model(tc.userA).Association("Follows").Find(&follows).Error
					assert.NoError(t, err)
					found := false
					for _, user := range follows {
						if user.ID == tc.userB.ID {
							found = true
							break
						}
					}
					assert.False(t, found, "userB should not be in userA's follows list")
				}
			}
		})
	}

	t.Run("Concurrent Unfollow", func(t *testing.T) {
		mockDB, err := gorm.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("failed to open mock database: %v", err)
		}
		defer mockDB.Close()

		store := &UserStore{db: mockDB}
		userA := createTestUser("userI")
		userB := createTestUser("userJ")

		mockDB.Model(userA).Association("Follows").Append(userB)

		var wg sync.WaitGroup
		concurrentCalls := 5

		for i := 0; i < concurrentCalls; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := store.Unfollow(userA, userB)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					t.Errorf("unexpected error in concurrent unfollow: %v", err)
				}
			}()
		}

		wg.Wait()

		var follows []model.User
		err = mockDB.Model(userA).Association("Follows").Find(&follows).Error
		assert.NoError(t, err)
		assert.Empty(t, follows, "follows list should be empty after concurrent unfollows")
	})
}

/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435


 */
func TestUpdate(t *testing.T) {

	tests := []struct {
		name    string
		user    *model.User
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successful Update",
			user: &model.User{
				Model: gorm.Model{
					ID:        1,
					UpdatedAt: time.Now(),
				},
				Username: "updated_user",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Bio:      "Updated bio",
				Image:    "updated.jpg",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Update with Duplicate Username",
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "existing_user",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(errors.New("Error 1062: Duplicate entry"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "Duplicate entry",
		},
		{
			name: "Update Non-Existent User",
			user: &model.User{
				Model: gorm.Model{ID: 999},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
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
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnError(errors.New("database connection lost"))
				mock.ExpectRollback()
			},
			wantErr: true,
			errMsg:  "database connection lost",
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
				t.Fatalf("Failed to open GORM DB: %v", err)
			}
			defer gormDB.Close()

			if tt.mockFn != nil {
				tt.mockFn(mock)
			}

			store := &UserStore{
				db: gormDB,
			}

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
				t.Errorf("Unfulfilled expectations: %s", err)
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
		name          string
		user          *model.User
		mockSetup     func(sqlmock.Sqlmock)
		expectedIDs   []uint
		expectedError error
	}{
		{
			name: "Successfully retrieve following user IDs",
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(uint(2)).
					AddRow(uint(3)).
					AddRow(uint(4))
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: nil,
		},
		{
			name: "User with no followings",
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name: "Database connection error",
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(uint(1)).
					WillReturnError(errors.New("database connection error"))
			},
			expectedIDs:   []uint{},
			expectedError: errors.New("database connection error"),
		},
		{
			name: "Invalid user ID",
			user: &model.User{
				Model: gorm.Model{ID: 999},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(uint(999)).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name: "Large number of followings",
			user: &model.User{
				Model: gorm.Model{ID: 1},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				for i := uint(1); i <= 1000; i++ {
					rows.AddRow(i)
				}
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			expectedIDs: func() []uint {
				ids := make([]uint, 1000)
				for i := uint(0); i < 1000; i++ {
					ids[i] = i + 1
				}
				return ids
			}(),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open GORM connection: %v", err)
			}
			defer gormDB.Close()

			tt.mockSetup(mock)

			store := &UserStore{db: gormDB}

			ids, err := store.GetFollowingUserIDs(tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedIDs, ids)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

