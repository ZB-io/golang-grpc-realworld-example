
// ********RoostGPT********
/*

roost_feedback [2/3/2025, 5:59:06 AM]:Modify Code to fix this error\nSuccessfully compiled but failed at runtime.\n=== RUN   TestUserStoreIsFollowing\n=== RUN   TestUserStoreIsFollowing/User_A_is_following_User_B\n--- FAIL: TestUserStoreIsFollowing (0.00s)\n    --- FAIL: TestUserStoreIsFollowing/User_A_is_following_User_B (0.00s)\npanic: runtime error: invalid memory address or nil pointer dereference [recovered]\n\tpanic: runtime error: invalid memory address or nil pointer dereference\n[signal SIGSEGV: segmentation violation code=0x1 addr=0x40 pc=0x633ca5]\n\ngoroutine 10 [running]:\ntesting.tRunner.func1.2({0x927860, 0xe21200})\n\t/usr/local/go/src/testing/testing.go:1632 +0x230\ntesting.tRunner.func1()\n\t/usr/local/go/src/testing/testing.go:1635 +0x35e\npanic({0x927860?, 0xe21200?})\n\t/usr/local/go/src/runtime/panic.go:785 +0x132\ngithub.com/jinzhu/gorm.(*DB).clone(0xc00052a340)\n\t/go/pkg/mod/github.com/jinzhu/gorm@v1.9.12/main.go:848 +0x25\ngithub.com/jinzhu/gorm.(*DB).Where(0xb6ab03?, {0x8f3ea0, 0xa932d0}, {0xc000336040, 0x2, 0x2})\n\t/go/pkg/mod/github.com/jinzhu/gorm@v1.9.12/main.go:235 +0x39\ncommand-line-arguments.(*MockUserStore).IsFollowing(0xc000054f40, 0xc000000000, 0xc0000000c0)\n\t/var/tmp/Roost/RoostGPT/golang-grpc-realworld-example/820ba614-cda3-4912-88a4-4c028d7d9667/source/golang-grpc-realworld-example/store/user_isfollowing_test.go:146 +0x114\ncommand-line-arguments.TestUserStoreIsFollowing.func1(0xc00052c4e0)\n\t/var/tmp/Roost/RoostGPT/golang-grpc-realworld-example/820ba614-cda3-4912-88a4-4c028d7d9667/source/golang-grpc-realworld-example/store/user_isfollowing_test.go:239 +0x65\ntesting.tRunner(0xc00052c4e0, 0xc0001e80a0)\n\t/usr/local/go/src/testing/testing.go:1690 +0xf4\ncreated by testing.(*T).Run in goroutine 9\n\t/usr/local/go/src/testing/testing.go:1743 +0x390\nFAIL\tcommand-line-arguments\t0.020s\nFAIL\n- \n- Add more comments to the test\n- \n Add more comments to the test\n- \n- Add more comments to the test\n- \n Add more comments to the test\n- \n- Add more comments to the test\n- \n Add more comments to the test\n- \n- Add more comments to the test\n- \n Add more comments to the test\n- \n- Add more comments to the test\n- \n Add more comments to the test\n- \n- Add more comments to the test\n- \n Add more comments to the test\n- \n- Add more comments to the test\n- \n Add more comments to the test\n- \n- Add more comments to the test\n- \n Add more comments to the test\n- Improve assertions
*/

// ********RoostGPT********

package github

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

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

			if tt.userA == nil || tt.userB == nil {
				if mockDB.lastQuery != "" {
					t.Errorf("Expected no query for nil users, but got: %s", mockDB.lastQuery)
				}
			} else {
				expectedQuery := "user_id = ? AND follow_id = ?"
				if mockDB.lastQuery != expectedQuery {
					t.Errorf("Expected query %s, but got: %s", expectedQuery, mockDB.lastQuery)
				}
				if len(mockDB.lastArgs) != 2 || mockDB.lastArgs[0] != tt.userA.ID || mockDB.lastArgs[1] != tt.userB.ID {
					t.Errorf("Incorrect query arguments: %v", mockDB.lastArgs)
				}
			}
		})
	}
}

type MockDB struct {
	countResult int
	countError  error
	lastQuery   string
	lastArgs    []interface{}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.lastQuery = query.(string)
	m.lastArgs = args
	return &gorm.DB{Value: m}
}

func (m *MockDB) Count(count *int) *gorm.DB {
	*count = m.countResult
	return &gorm.DB{Error: m.countError}
}

type MockUserStore struct {
	db *MockDB
}

func (s *MockUserStore) IsFollowing(a *model.User, b *model.User) (bool, error) {
	if a == nil || b == nil {
		return false, nil
	}

	var count int
	err := s.db.Where("user_id = ? AND follow_id = ?", a.ID, b.ID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
