package github

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
)









/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore 

 */
func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
		want *UserStore
	}{
		{
			name: "Create UserStore with valid DB",
			db:   &gorm.DB{},
			want: &UserStore{db: &gorm.DB{}},
		},
		{
			name: "Create UserStore with nil DB",
			db:   nil,
			want: &UserStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)

			if got == nil {
				t.Errorf("NewUserStore() returned nil")
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserStore() = %v, want %v", got, tt.want)
			}

			if got.db != tt.db {
				t.Errorf("NewUserStore().db = %v, want %v", got.db, tt.db)
			}

			if reflect.TypeOf(got) != reflect.TypeOf(&UserStore{}) {
				t.Errorf("NewUserStore() returned type %v, want *UserStore", reflect.TypeOf(got))
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
		name        string
		userA       *model.User
		userB       *model.User
		countResult int
		countError  error
		want        bool
		wantErr     bool
	}{
		{
			name:        "User A is following User B",
			userA:       &model.User{Model: gorm.Model{ID: 1}},
			userB:       &model.User{Model: gorm.Model{ID: 2}},
			countResult: 1,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "User A is not following User B",
			userA:       &model.User{Model: gorm.Model{ID: 1}},
			userB:       &model.User{Model: gorm.Model{ID: 2}},
			countResult: 0,
			want:        false,
			wantErr:     false,
		},
		{
			name:    "Both users are nil",
			userA:   nil,
			userB:   nil,
			want:    false,
			wantErr: false,
		},
		{
			name:    "User A is nil",
			userA:   nil,
			userB:   &model.User{Model: gorm.Model{ID: 2}},
			want:    false,
			wantErr: false,
		},
		{
			name:    "User B is nil",
			userA:   &model.User{Model: gorm.Model{ID: 1}},
			userB:   nil,
			want:    false,
			wantErr: false,
		},
		{
			name:       "Database error occurs",
			userA:      &model.User{Model: gorm.Model{ID: 1}},
			userB:      &model.User{Model: gorm.Model{ID: 2}},
			countError: errors.New("database error"),
			want:       false,
			wantErr:    true,
		},
		{
			name:        "User is following themselves",
			userA:       &model.User{Model: gorm.Model{ID: 1}},
			userB:       &model.User{Model: gorm.Model{ID: 1}},
			countResult: 1,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "Users exist but have no relationship",
			userA:       &model.User{Model: gorm.Model{ID: 1}},
			userB:       &model.User{Model: gorm.Model{ID: 2}},
			countResult: 0,
			want:        false,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{
				countResult: tt.countResult,
				countError:  tt.countError,
			}

			s := &MockUserStore{
				db: mockDB,
			}

			got, err := s.IsFollowing(tt.userA, tt.userB)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.IsFollowing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserStore.IsFollowing() = %v, want %v", got, tt.want)
			}
		})
	}
}

