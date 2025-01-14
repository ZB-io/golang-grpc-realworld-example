package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"time"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"database/sql"
)






type mockDB struct {
	gorm.DB
	updateFunc func(interface{}) *gorm.DB
}
type mockDB struct {
	count int
	err   error
}


/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore 

 */
func BenchmarkNewUserStore(b *testing.B) {
	db := &gorm.DB{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		NewUserStore(db)
	}
}

func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
		want *UserStore
	}{
		{
			name: "Create UserStore with valid gorm.DB instance",
			db:   &gorm.DB{},
			want: &UserStore{db: &gorm.DB{}},
		},
		{
			name: "Create UserStore with nil gorm.DB instance",
			db:   nil,
			want: &UserStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserStoreImmutability(t *testing.T) {
	db := &gorm.DB{}
	store1 := NewUserStore(db)
	store2 := NewUserStore(db)

	if store1 == store2 {
		t.Error("NewUserStore() returned the same instance for multiple calls")
	}
}

func TestNewUserStoreWithConfiguredDB(t *testing.T) {
	db := &gorm.DB{}

	store := NewUserStore(db)

	if !reflect.DeepEqual(store.db, db) {
		t.Error("NewUserStore() did not preserve the configuration of the input gorm.DB instance")
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
				return db.InstantSet("gorm:get_first_result", &model.User{
					Model:    gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Username: "testuser",
					Email:    "testuser@example.com",
				})
			},
			want: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "testuser@example.com",
			},
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
			username: "maxlengthusername",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				return db.InstantSet("gorm:get_first_result", &model.User{
					Model:    gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Username: "maxlengthusername",
					Email:    "maxlength@example.com",
				})
			},
			want: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "maxlengthusername",
				Email:    "maxlength@example.com",
			},
			wantErr: nil,
		},
		{
			name:     "Attempt retrieval with an empty username",
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
			name:     "Handle case sensitivity in username lookup",
			username: "TestUser",
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.Error = gorm.ErrRecordNotFound
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
func (m *mockDB) Model(value interface{}) *gorm.DB {
	return &gorm.DB{Value: value}
}

func TestUserStoreUpdate(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		mockDB  func() *gorm.DB
		wantErr bool
	}{
		{
			name: "Successful User Update",
			user: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "updateduser",
				Email:    "updated@example.com",
			},
			mockDB: func() *gorm.DB {
				mock := &mockDB{
					updateFunc: func(interface{}) *gorm.DB {
						return &gorm.DB{Error: nil}
					},
				}
				return &mock.DB
			},
			wantErr: false,
		},
		{
			name: "Database Error During Update",
			user: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "erroruser",
			},
			mockDB: func() *gorm.DB {
				mock := &mockDB{
					updateFunc: func(interface{}) *gorm.DB {
						return &gorm.DB{Error: errors.New("database error")}
					},
				}
				return &mock.DB
			},
			wantErr: true,
		},
		{
			name: "Update with Empty User Model",
			user: &model.User{},
			mockDB: func() *gorm.DB {
				mock := &mockDB{
					updateFunc: func(interface{}) *gorm.DB {
						return &gorm.DB{Error: nil}
					},
				}
				return &mock.DB
			},
			wantErr: false,
		},
		{
			name: "Update User with Changed Primary Key",
			user: &model.User{
				Model:    gorm.Model{ID: 3},
				Username: "changediduser",
			},
			mockDB: func() *gorm.DB {
				mock := &mockDB{
					updateFunc: func(interface{}) *gorm.DB {
						return &gorm.DB{Error: nil}
					},
				}
				return &mock.DB
			},
			wantErr: false,
		},
		{
			name: "Update with Invalid User Data",
			user: &model.User{
				Model:    gorm.Model{ID: 4},
				Username: "invaliduser",
				Email:    "invalid",
			},
			mockDB: func() *gorm.DB {
				mock := &mockDB{
					updateFunc: func(interface{}) *gorm.DB {
						return &gorm.DB{Error: errors.New("constraint violation")}
					},
				}
				return &mock.DB
			},
			wantErr: true,
		},
		{
			name: "Partial User Update",
			user: &model.User{
				Model: gorm.Model{ID: 5},
				Bio:   "Updated bio",
			},
			mockDB: func() *gorm.DB {
				mock := &mockDB{
					updateFunc: func(interface{}) *gorm.DB {
						return &gorm.DB{Error: nil}
					},
				}
				return &mock.DB
			},
			wantErr: false,
		},
		{
			name: "Update User with Associated Data",
			user: &model.User{
				Model:            gorm.Model{ID: 6},
				Username:         "associateduser",
				Follows:          []model.User{{Model: gorm.Model{ID: 7}}},
				FavoriteArticles: []model.Article{{Model: gorm.Model{ID: 8}}},
			},
			mockDB: func() *gorm.DB {
				mock := &mockDB{
					updateFunc: func(interface{}) *gorm.DB {
						return &gorm.DB{Error: nil}
					},
				}
				return &mock.DB
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := tt.mockDB()
			store := &UserStore{db: mockDB}

			err := store.Update(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (m *mockDB) Update(attrs ...interface{}) *gorm.DB {
	return m.updateFunc(attrs[0])
}


/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c

FUNCTION_DEF=func (s *UserStore) IsFollowing(a *model.User, b *model.User) (bool, error) 

 */
func (m *mockDB) Count(value interface{}) *gorm.DB {
	*(value.(*int)) = m.count
	return &gorm.DB{Error: m.err}
}

func (m *mockDB) Table(name string) *gorm.DB {
	return &gorm.DB{Value: m}
}

func TestUserStoreIsFollowing(t *testing.T) {
	tests := []struct {
		name     string
		a        *model.User
		b        *model.User
		mockDB   *mockDB
		expected bool
		wantErr  bool
	}{
		{
			name:     "User A is following User B",
			a:        &model.User{Model: gorm.Model{ID: 1}},
			b:        &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   &mockDB{count: 1, err: nil},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "User A is not following User B",
			a:        &model.User{Model: gorm.Model{ID: 1}},
			b:        &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   &mockDB{count: 0, err: nil},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Null user arguments",
			a:        nil,
			b:        nil,
			mockDB:   &mockDB{},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Database error",
			a:        &model.User{Model: gorm.Model{ID: 1}},
			b:        &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   &mockDB{err: errors.New("database error")},
			expected: false,
			wantErr:  true,
		},
		{
			name:     "User following themselves",
			a:        &model.User{Model: gorm.Model{ID: 1}},
			b:        &model.User{Model: gorm.Model{ID: 1}},
			mockDB:   &mockDB{count: 1, err: nil},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Users with no ID (unsaved users)",
			a:        &model.User{},
			b:        &model.User{},
			mockDB:   &mockDB{count: 0, err: nil},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "Large number of follows",
			a:        &model.User{Model: gorm.Model{ID: 1}},
			b:        &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   &mockDB{count: 1000, err: nil},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "Deleted user relationship",
			a:        &model.User{Model: gorm.Model{ID: 1}},
			b:        &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   &mockDB{count: 0, err: nil},
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserStore{
				db: &gorm.DB{Value: tt.mockDB},
			}

			got, err := s.IsFollowing(tt.a, tt.b)
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

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return &gorm.DB{Value: m}
}


/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7

FUNCTION_DEF=func (s *UserStore) GetFollowingUserIDs(m *model.User) ([]uint, error) 

 */
func (m *mockDB) Rows() (*sql.Rows, error) {
	if m.err != nil {
		return nil, m.err
	}

	return nil, nil
}

func (m *mockDB) Select(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}

func (m *mockDB) Table(name string) *gorm.DB {
	return m.DB
}

func TestUserStoreGetFollowingUserIDs(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		mockDB  func() *gorm.DB
		want    []uint
		wantErr bool
	}{
		{
			name: "Successfully retrieve following user IDs",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			want:    []uint{2, 3, 4},
			wantErr: false,
		},
		{
			name: "User with no followers",
			user: &model.User{Model: gorm.Model{ID: 2}},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			want:    []uint{},
			wantErr: false,
		},
		{
			name: "Database error",
			user: &model.User{Model: gorm.Model{ID: 3}},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(errors.New("database error"))
				return db
			},
			want:    []uint{},
			wantErr: true,
		},
		{
			name: "User with large number of followers",
			user: &model.User{Model: gorm.Model{ID: 4}},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			want:    generateLargeIDSlice(1000),
			wantErr: false,
		},
		{
			name: "Invalid user ID",
			user: &model.User{Model: gorm.Model{ID: 9999}},
			mockDB: func() *gorm.DB {
				db := &gorm.DB{}
				db.AddError(nil)
				return db
			},
			want:    []uint{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserStore{
				db: tt.mockDB(),
			}

			got, err := s.GetFollowingUserIDs(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.GetFollowingUserIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStore.GetFollowingUserIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (m *mockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB
}

func generateLargeIDSlice(n int) []uint {
	ids := make([]uint, n)
	for i := 0; i < n; i++ {
		ids[i] = uint(i + 1)
	}
	return ids
}

