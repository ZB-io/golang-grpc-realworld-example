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
)

/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920


 */
func TestCreate(t *testing.T) {

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
				Bio:      "",
				Image:    "",
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
			dbError: errors.New("database connection failed"),
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

			gdb, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("Failed to open gorm DB: %v", err)
			}
			defer gdb.Close()

			store := &UserStore{
				db: gdb,
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

	tests := []struct {
		name    string
		userA   *model.User
		userB   *model.User
		dbSetup func(*gorm.DB)
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
			dbSetup: func(db *gorm.DB) {

			},
			wantErr: false,
		},
		{
			name:  "Nil user A",
			userA: nil,
			userB: &model.User{
				Model: gorm.Model{ID: 2},
			},
			wantErr: true,
			errMsg:  "invalid user parameters",
		},
		{
			name: "Nil user B",
			userA: &model.User{
				Model: gorm.Model{ID: 1},
			},
			userB:   nil,
			wantErr: true,
			errMsg:  "invalid user parameters",
		},
		{
			name: "Self follow attempt",
			userA: &model.User{
				Model: gorm.Model{ID: 1},
			},
			userB: &model.User{
				Model: gorm.Model{ID: 1},
			},
			wantErr: true,
			errMsg:  "self-follow not allowed",
		},
		{
			name: "Already following user",
			userA: &model.User{
				Model: gorm.Model{ID: 1},
			},
			userB: &model.User{
				Model: gorm.Model{ID: 2},
			},
			dbSetup: func(db *gorm.DB) {

			},
			wantErr: false,
		},
		{
			name: "Database error",
			userA: &model.User{
				Model: gorm.Model{ID: 1},
			},
			userB: &model.User{
				Model: gorm.Model{ID: 2},
			},
			dbSetup: func(db *gorm.DB) {

			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db := &gorm.DB{}

			if tt.dbSetup != nil {
				tt.dbSetup(db)
			}

			store := &UserStore{
				db: db,
			}

			err := store.Follow(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("Follow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Follow() error message = %v, want %v", err.Error(), tt.errMsg)
			}

			if !tt.wantErr {
				var follows []model.User
				if err := db.Model(tt.userA).Association("Follows").Find(&follows).Error; err != nil {
					t.Errorf("Failed to verify follow relationship: %v", err)
				}

				found := false
				for _, f := range follows {
					if f.ID == tt.userB.ID {
						found = true
						break
					}
				}

				if !found {
					t.Error("Follow relationship not found in database")
				}
			}
		})
	}
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
			name:  "Successful user retrieval",
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
			mockSetup: func(mock *MockDB) {
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
			mockSetup: func(mock *MockDB) {
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
			mockSetup: func(mock *MockDB) {
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

			mockDB := new(MockDB)
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
		name        string
		id          uint
		setupMock   func(*MockDB)
		expectUser  *model.User
		expectError error
	}{
		{
			name: "Successfully retrieve user by valid ID",
			id:   1,
			setupMock: func(m *MockDB) {
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
			setupMock: func(m *MockDB) {
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
			setupMock: func(m *MockDB) {
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
			setupMock: func(m *MockDB) {
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

			mockDB := new(MockDB)
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
func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	called := m.Called(out, where)
	return called.Get(0).(*gorm.DB)
}

func TestGetByUsername(t *testing.T) {

	tests := []struct {
		name          string
		username      string
		mockSetup     func(*MockDB)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:     "Successful user retrieval",
			username: "validuser",
			mockSetup: func(mock *MockDB) {
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
			mockSetup: func(mock *MockDB) {
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
			mockSetup: func(mock *MockDB) {
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
			mockSetup: func(mock *MockDB) {
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

			mockDB := new(MockDB)
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

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	called := m.Called(query, args)
	return called.Get(0).(*gorm.DB)
}

