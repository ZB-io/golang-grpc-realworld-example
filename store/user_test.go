package store

import (
	"sync"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"math"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/DATA-DOG/go-sqlmock"
)






type MockDB struct {
	error        error
	modelCalled  bool
	updateCalled bool
}


/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore 

 */
func (m *mockLogger) Print(v ...interface{}) {}

func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name     string
		db       *gorm.DB
		wantNil  bool
		wantSame bool
	}{
		{
			name:     "Valid DB Connection",
			db:       &gorm.DB{},
			wantNil:  false,
			wantSame: true,
		},
		{
			name:     "Nil DB Connection",
			db:       nil,
			wantNil:  false,
			wantSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)

			if (got == nil) != tt.wantNil {
				t.Errorf("NewUserStore() returned nil: %v, want nil: %v", got == nil, tt.wantNil)
			}

			if !tt.wantNil {
				if (got.db == tt.db) != tt.wantSame {
					t.Errorf("NewUserStore().db same as input: %v, want same: %v", got.db == tt.db, tt.wantSame)
				}
			}
		})
	}
}

func TestNewUserStoreConcurrency(t *testing.T) {
	db := &gorm.DB{}
	numGoroutines := 100
	var wg sync.WaitGroup
	stores := make([]*UserStore, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			stores[index] = NewUserStore(db)
		}(i)
	}

	wg.Wait()

	for i := 0; i < numGoroutines; i++ {
		if stores[i] == nil {
			t.Errorf("NewUserStore() returned nil in goroutine %d", i)
		}
		if stores[i].db != db {
			t.Errorf("NewUserStore() returned instance with incorrect db reference in goroutine %d", i)
		}
	}
}

func TestNewUserStoreImmutability(t *testing.T) {
	db := &gorm.DB{}
	store1 := NewUserStore(db)
	store2 := NewUserStore(db)

	if store1 == store2 {
		t.Error("NewUserStore() returned the same instance for multiple calls")
	}

	if store1.db != store2.db {
		t.Error("NewUserStore() returned instances with different db references")
	}
}

func TestNewUserStoreWithCustomConfig(t *testing.T) {

	db := &gorm.DB{}

	customLogger := &mockLogger{}
	db.SetLogger(customLogger)

	store := NewUserStore(db)

	if store.db != db {
		t.Error("NewUserStore() did not preserve the input db")
	}

	store.db.LogMode(true)

}


/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1

FUNCTION_DEF=func (s *UserStore) GetByID(id uint) (*model.User, error) 

 */
func TestUserStoreGetById(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		mockDB  func() *gorm.DB
		want    *model.User
		wantErr error
	}{
		{
			name: "Successful Retrieval of Existing User",
			id:   1,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.AutoMigrate(&model.User{})
				user := model.User{Model: gorm.Model{ID: 1}, Username: "testuser", Email: "test@example.com"}
				db.Create(&user)
				return db
			},
			want: &model.User{Model: gorm.Model{ID: 1}, Username: "testuser", Email: "test@example.com"},
		},
		{
			name: "Attempt to Retrieve Non-existent User",
			id:   999,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.AutoMigrate(&model.User{})
				return db
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			id:   1,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.Close()
				return db
			},
			want:    nil,
			wantErr: errors.New("sql: database is closed"),
		},
		{
			name: "Retrieval with Zero ID",
			id:   0,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.AutoMigrate(&model.User{})
				return db
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Retrieval of Soft-Deleted User",
			id:   2,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.AutoMigrate(&model.User{})
				user := model.User{Model: gorm.Model{ID: 2}, Username: "deleteduser", Email: "deleted@example.com"}
				db.Create(&user)
				db.Delete(&user)
				return db
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Retrieval with Maximum uint Value",
			id:   math.MaxUint32,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.AutoMigrate(&model.User{})
				return db
			},
			want:    nil,
			wantErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserStore{
				db: tt.mockDB(),
			}

			got, err := s.GetByID(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435

FUNCTION_DEF=func (s *UserStore) Update(m *model.User) error 

 */
func (m *MockDB) Error() error {
	return m.error
}

func (m *MockDB) Model(value interface{}) *MockDB {
	m.modelCalled = true
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
				Model: gorm.Model{ID: 2},
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Update with Database Error",
			user: &model.User{
				Model:    gorm.Model{ID: 3},
				Username: "erroruser",
			},
			dbError: errors.New("database error"),
			wantErr: true,
		},
		{
			name:    "Update with Nil User",
			user:    nil,
			dbError: nil,
			wantErr: true,
		},
		{
			name: "Update Affecting Multiple Fields",
			user: &model.User{
				Model:    gorm.Model{ID: 4},
				Username: "multiupdate",
				Email:    "multi@example.com",
				Bio:      "Updated bio",
			},
			dbError: nil,
			wantErr: false,
		},
		{
			name: "Update with Unique Constraint Violation",
			user: &model.User{
				Model: gorm.Model{ID: 5},
				Email: "existing@example.com",
			},
			dbError: errors.New("unique constraint violation"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &MockDB{
				error: tt.dbError,
			}

			s := &MockUserStore{
				db: mockDB,
			}

			err := s.Update(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.user != nil {
				if !mockDB.modelCalled {
					t.Errorf("Model() was not called")
				}
				if !mockDB.updateCalled {
					t.Errorf("Update() was not called")
				}
			}
		})
	}
}

func (m *MockDB) Update(attrs ...interface{}) *MockDB {
	m.updateCalled = true
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
		setupMock   func(*gorm.DB)
		follower    *model.User
		followed    *model.User
		expectedErr error
	}{
		{
			name: "Successful Follow Operation",
			setupMock: func(db *gorm.DB) {
				db.AddError(nil)
			},
			follower:    &model.User{Model: gorm.Model{ID: 1}},
			followed:    &model.User{Model: gorm.Model{ID: 2}},
			expectedErr: nil,
		},
		{
			name: "Follow User Already Being Followed",
			setupMock: func(db *gorm.DB) {
				db.AddError(nil)
			},
			follower:    &model.User{Model: gorm.Model{ID: 1}},
			followed:    &model.User{Model: gorm.Model{ID: 2}},
			expectedErr: nil,
		},
		{
			name: "Self-Follow Attempt",
			setupMock: func(db *gorm.DB) {
				db.AddError(nil)
			},
			follower:    &model.User{Model: gorm.Model{ID: 1}},
			followed:    &model.User{Model: gorm.Model{ID: 1}},
			expectedErr: nil,
		},
		{
			name: "Follow Non-Existent User",
			setupMock: func(db *gorm.DB) {
				db.AddError(gorm.ErrRecordNotFound)
			},
			follower:    &model.User{Model: gorm.Model{ID: 1}},
			followed:    &model.User{Model: gorm.Model{ID: 999}},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			setupMock: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
			follower:    &model.User{Model: gorm.Model{ID: 1}},
			followed:    &model.User{Model: gorm.Model{ID: 2}},
			expectedErr: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &gorm.DB{}
			tt.setupMock(mockDB)

			store := &UserStore{db: mockDB}

			err := store.Follow(tt.follower, tt.followed)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Follow() error = %v, expectedErr %v", err, tt.expectedErr)
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
	tests := []struct {
		name     string
		userA    *model.User
		userB    *model.User
		dbSetup  func(*gorm.DB)
		expected bool
		err      error
	}{
		{
			name:  "User A is following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			dbSetup: func(db *gorm.DB) {
				db.Table("follows").Create(map[string]interface{}{
					"from_user_id": 1,
					"to_user_id":   2,
				})
			},
			expected: true,
			err:      nil,
		},
		{
			name:  "User A is not following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			dbSetup: func(db *gorm.DB) {

			},
			expected: false,
			err:      nil,
		},
		{
			name:  "User A is following User B among many other followings",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			dbSetup: func(db *gorm.DB) {
				db.Table("follows").Create(map[string]interface{}{
					"from_user_id": 1,
					"to_user_id":   2,
				})
				db.Table("follows").Create(map[string]interface{}{
					"from_user_id": 1,
					"to_user_id":   3,
				})
				db.Table("follows").Create(map[string]interface{}{
					"from_user_id": 2,
					"to_user_id":   3,
				})
			},
			expected: true,
			err:      nil,
		},
		{
			name:  "Database error occurs",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			dbSetup: func(db *gorm.DB) {

				db.AddError(errors.New("database error"))
			},
			expected: false,
			err:      errors.New("database error"),
		},
		{
			name:     "User A parameter is nil",
			userA:    nil,
			userB:    &model.User{Model: gorm.Model{ID: 2}},
			dbSetup:  func(db *gorm.DB) {},
			expected: false,
			err:      nil,
		},
		{
			name:     "User B parameter is nil",
			userA:    &model.User{Model: gorm.Model{ID: 1}},
			userB:    nil,
			dbSetup:  func(db *gorm.DB) {},
			expected: false,
			err:      nil,
		},
		{
			name:     "Both user parameters are nil",
			userA:    nil,
			userB:    nil,
			dbSetup:  func(db *gorm.DB) {},
			expected: false,
			err:      nil,
		},
		{
			name:  "User is following themselves",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 1}},
			dbSetup: func(db *gorm.DB) {
				db.Table("follows").Create(map[string]interface{}{
					"from_user_id": 1,
					"to_user_id":   1,
				})
			},
			expected: true,
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, _ := gorm.Open("sqlite3", ":memory:")
			defer db.Close()

			tt.dbSetup(db)

			store := &UserStore{db: db}

			result, err := store.IsFollowing(tt.userA, tt.userB)

			if result != tt.expected {
				t.Errorf("Expected result %v, but got %v", tt.expected, result)
			}

			if (err != nil && tt.err == nil) || (err == nil && tt.err != nil) || (err != nil && tt.err != nil && err.Error() != tt.err.Error()) {
				t.Errorf("Expected error %v, but got %v", tt.err, err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55

FUNCTION_DEF=func (s *UserStore) Unfollow(a *model.User, b *model.User) error 

 */
func TestUserStoreUnfollow(t *testing.T) {
	tests := []struct {
		name    string
		userA   *model.User
		userB   *model.User
		mockDB  func() *gorm.DB
		wantErr bool
		errMsg  string
	}{
		{
			name:  "Successfully Unfollow a User",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			wantErr: false,
		},
		{
			name:  "Attempt to Unfollow a User Not Currently Followed",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			wantErr: false,
		},
		{
			name:  "Handle Database Error During Unfollow Operation",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(errors.New("database error"))
				return db
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:  "Unfollow with Nil User Arguments (userA)",
			userA: nil,
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDB: func() *gorm.DB {
				return &gorm.DB{}
			},
			wantErr: true,
			errMsg:  "invalid user argument",
		},
		{
			name:  "Unfollow with Nil User Arguments (userB)",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: nil,
			mockDB: func() *gorm.DB {
				return &gorm.DB{}
			},
			wantErr: true,
			errMsg:  "invalid user argument",
		},
		{
			name:  "Unfollow Self",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func() *gorm.DB {
				return &gorm.DB{}
			},
			wantErr: true,
			errMsg:  "cannot unfollow self",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserStore{
				db: tt.mockDB(),
			}

			err := s.Unfollow(tt.userA, tt.userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Unfollow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("UserStore.Unfollow() error message = %v, want %v", err.Error(), tt.errMsg)
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
			name: "Successful retrieval of following user IDs",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3).
					AddRow(4)
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    []uint{2, 3, 4},
			wantErr: false,
		},
		{
			name: "User with no followers",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    []uint{},
			wantErr: false,
		},
		{
			name: "Database error handling",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			want:    []uint{},
			wantErr: true,
		},
		{
			name: "Large number of followers",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				for i := 2; i <= 1001; i++ {
					rows.AddRow(uint(i))
				}
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: func() []uint {
				ids := make([]uint, 1000)
				for i := range ids {
					ids[i] = uint(i + 2)
				}
				return ids
			}(),
			wantErr: false,
		},
		{
			name: "Invalid user ID",
			user: &model.User{Model: gorm.Model{ID: 999}},
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(999).
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

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a gorm database", err)
			}

			tt.mockDB(mock)

			s := &UserStore{db: gormDB}

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

