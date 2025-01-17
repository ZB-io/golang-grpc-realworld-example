package store

import (
	"reflect"
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"errors"
	"fmt"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-sqlmock"
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
			name: "Create UserStore with valid gorm.DB",
			db:   &gorm.DB{},
			want: &UserStore{db: &gorm.DB{}},
		},
		{
			name: "Create UserStore with nil gorm.DB",
			db:   nil,
			want: &UserStore{db: nil},
		},
		{
			name: "Create UserStore with configured gorm.DB",
			db: &gorm.DB{
				Error:        nil,
				RowsAffected: 10,
				Value:        "test",
			},
			want: &UserStore{
				db: &gorm.DB{
					Error:        nil,
					RowsAffected: 10,
					Value:        "test",
				},
			},
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

func TestNewUserStore_Concurrency(t *testing.T) {
	iterations := 100
	ch := make(chan *UserStore, iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			db := &gorm.DB{}
			ch <- NewUserStore(db)
		}()
	}

	for i := 0; i < iterations; i++ {
		select {
		case store := <-ch:
			if store == nil {
				t.Error("NewUserStore() returned nil in concurrent execution")
			}
		case <-time.After(time.Second):
			t.Error("NewUserStore() timed out in concurrent execution")
		}
	}
}

func TestNewUserStore_Immutability(t *testing.T) {
	db := &gorm.DB{}
	store1 := NewUserStore(db)
	store2 := NewUserStore(db)

	if store1 == store2 {
		t.Error("NewUserStore() should return different instances")
	}

	if store1.db != store2.db {
		t.Error("NewUserStore() instances should share the same db")
	}
}

func TestNewUserStore_Performance(t *testing.T) {
	db := &gorm.DB{}
	iterations := 1000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		NewUserStore(db)
	}

	duration := time.Since(start)
	avgDuration := duration / time.Duration(iterations)

	if avgDuration > time.Millisecond {
		t.Errorf("NewUserStore() average duration %v exceeds 1ms threshold", avgDuration)
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
		mockDB   func() *gorm.DB
		expected *model.User
		wantErr  bool
		err      error
	}{
		{
			name: "Successfully retrieve a user by ID",
			id:   1,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				user := &model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password",
				}
				db.Create(user)
				return db
			},
			expected: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "Attempt to retrieve a non-existent user",
			id:   999,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				return db
			},
			expected: nil,
			wantErr:  true,
			err:      gorm.ErrRecordNotFound,
		},
		{
			name: "Database connection error",
			id:   1,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				db.Close()
				return db
			},
			expected: nil,
			wantErr:  true,
			err:      errors.New("sql: database is closed"),
		},
		{
			name: "Retrieve user with minimum fields populated",
			id:   2,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				user := &model.User{
					Model:    gorm.Model{ID: 2},
					Username: "minuser",
					Email:    "min@example.com",
					Password: "minpass",
				}
				db.Create(user)
				return db
			},
			expected: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "minuser",
				Email:    "min@example.com",
				Password: "minpass",
			},
			wantErr: false,
		},
		{
			name: "Retrieve user with all fields populated",
			id:   3,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				user := &model.User{
					Model:    gorm.Model{ID: 3},
					Username: "fulluser",
					Email:    "full@example.com",
					Password: "fullpass",
					Bio:      "Full bio",
					Image:    "full.jpg",
					Follows: []model.User{
						{Model: gorm.Model{ID: 4}, Username: "follower"},
					},
					FavoriteArticles: []model.Article{
						{Model: gorm.Model{ID: 1}, Title: "Favorite Article"},
					},
				}
				db.Create(user)
				return db
			},
			expected: &model.User{
				Model:    gorm.Model{ID: 3},
				Username: "fulluser",
				Email:    "full@example.com",
				Password: "fullpass",
				Bio:      "Full bio",
				Image:    "full.jpg",
				Follows: []model.User{
					{Model: gorm.Model{ID: 4}, Username: "follower"},
				},
				FavoriteArticles: []model.Article{
					{Model: gorm.Model{ID: 1}, Title: "Favorite Article"},
				},
			},
			wantErr: false,
		},
		{
			name: "Performance with a large number of users",
			id:   50000,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				for i := 1; i <= 100000; i++ {
					user := &model.User{
						Model:    gorm.Model{ID: uint(i)},
						Username: fmt.Sprintf("user%d", i),
						Email:    fmt.Sprintf("user%d@example.com", i),
						Password: fmt.Sprintf("pass%d", i),
					}
					db.Create(user)
				}
				return db
			},
			expected: &model.User{
				Model:    gorm.Model{ID: 50000},
				Username: "user50000",
				Email:    "user50000@example.com",
				Password: "pass50000",
			},
			wantErr: false,
		},
		{
			name: "Behavior with zero ID",
			id:   0,
			mockDB: func() *gorm.DB {
				db, _ := gorm.Open("sqlite3", ":memory:")
				return db
			},
			expected: nil,
			wantErr:  true,
			err:      gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.mockDB()
			s := &UserStore{db: db}

			start := time.Now()
			user, err := s.GetByID(tt.id)
			duration := time.Since(start)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.err != nil {
					assert.EqualError(t, err, tt.err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}

			if tt.name == "Performance with a large number of users" {
				assert.Less(t, duration, 1*time.Second, "GetByID took too long for large dataset")
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435

FUNCTION_DEF=func (s *UserStore) Update(m *model.User) error 

 */
func TestUserStoreUpdate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*gorm.DB)
		input   *model.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successfully Update User Information",
			setup: func(db *gorm.DB) {
				db.AddError(nil)
			},
			input: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "updateduser",
				Email:    "updated@example.com",
				Bio:      "Updated bio",
			},
			wantErr: false,
		},
		{
			name: "Attempt to Update Non-Existent User",
			setup: func(db *gorm.DB) {
				db.AddError(gorm.ErrRecordNotFound)
			},
			input: &model.User{
				Model: gorm.Model{ID: 999},
			},
			wantErr: true,
			errMsg:  "record not found",
		},
		{
			name: "Update User with Invalid Data",
			setup: func(db *gorm.DB) {
				db.AddError(errors.New("validation error"))
			},
			input: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "",
			},
			wantErr: true,
			errMsg:  "validation error",
		},
		{
			name: "Update User with Duplicate Unique Fields",
			setup: func(db *gorm.DB) {
				db.AddError(errors.New("unique constraint violation"))
			},
			input: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "existinguser",
			},
			wantErr: true,
			errMsg:  "unique constraint violation",
		},
		{
			name: "Partial Update of User Information",
			setup: func(db *gorm.DB) {
				db.AddError(nil)
			},
			input: &model.User{
				Model: gorm.Model{ID: 1},
				Bio:   "Updated bio",
				Image: "new-image.jpg",
			},
			wantErr: false,
		},
		{
			name: "Update User with No Changes",
			setup: func(db *gorm.DB) {
				db.AddError(nil)
			},
			input: &model.User{
				Model: gorm.Model{ID: 1},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &gorm.DB{}
			tt.setup(mockDB)

			us := &UserStore{db: mockDB}

			err := us.Update(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.errMsg)
			}

			if !reflect.DeepEqual(mockDB.Value, tt.input) {
				t.Errorf("Update() did not update with correct user data")
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Follow_48fdf1257b
ROOST_METHOD_SIG_HASH=Follow_8217e61c06

FUNCTION_DEF=func (s *UserStore) Follow(a *model.User, b *model.User) error 

 */
func TestUserStoreFollow(t *testing.T) {
	type mockDB struct {
		*gorm.DB
	}

	type args struct {
		a *model.User
		b *model.User
	}

	tests := []struct {
		name    string
		db      *mockDB
		args    args
		wantErr bool
	}{
		{
			name: "Successful Follow Operation",
			db: &mockDB{
				DB: &gorm.DB{},
			},
			args: args{
				a: &model.User{Model: gorm.Model{ID: 1}},
				b: &model.User{Model: gorm.Model{ID: 2}},
			},
			wantErr: false,
		},
		{
			name: "Follow User Already Being Followed",
			db: &mockDB{
				DB: &gorm.DB{},
			},
			args: args{
				a: &model.User{Model: gorm.Model{ID: 1}, Follows: []model.User{{Model: gorm.Model{ID: 2}}}},
				b: &model.User{Model: gorm.Model{ID: 2}},
			},
			wantErr: false,
		},
		{
			name: "Follow Non-Existent User",
			db: &mockDB{
				DB: &gorm.DB{},
			},
			args: args{
				a: &model.User{Model: gorm.Model{ID: 1}},
				b: &model.User{Model: gorm.Model{ID: 999}},
			},
			wantErr: true,
		},
		{
			name: "Follow Self",
			db: &mockDB{
				DB: &gorm.DB{},
			},
			args: args{
				a: &model.User{Model: gorm.Model{ID: 1}},
				b: &model.User{Model: gorm.Model{ID: 1}},
			},
			wantErr: false,
		},
		{
			name: "Follow with Nil User",
			db: &mockDB{
				DB: &gorm.DB{},
			},
			args: args{
				a: &model.User{Model: gorm.Model{ID: 1}},
				b: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserStore{
				db: tt.db.DB,
			}

			tt.db.DB = tt.db.DB.InstantSet("gorm:association_mock", func(s *gorm.Association) error {
				if tt.args.b == nil {
					return errors.New("cannot follow nil user")
				}
				if tt.args.b.ID == 999 {
					return errors.New("user not found")
				}
				return nil
			})

			err := s.Follow(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Follow() error = %v, wantErr %v", err, tt.wantErr)
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
		name          string
		user          *model.User
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedIDs   []uint
		expectedError error
	}{
		{
			name: "Successful retrieval of following user IDs",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
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
			name: "User with no followers",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name: "Database error handling",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			expectedIDs:   []uint{},
			expectedError: errors.New("database error"),
		},
		{
			name: "Large number of followed users",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				for i := 2; i <= 1001; i++ {
					rows.AddRow(uint(i))
				}
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   generateSequence(2, 1001),
			expectedError: nil,
		},
		{
			name: "Handling of deleted users",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3)
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedIDs:   []uint{2, 3},
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
				t.Fatalf("Failed to open GORM DB: %v", err)
			}

			tt.mockSetup(mock)

			store := &UserStore{db: gormDB}
			gotIDs, gotErr := store.GetFollowingUserIDs(tt.user)

			if !reflect.DeepEqual(gotIDs, tt.expectedIDs) {
				t.Errorf("GetFollowingUserIDs() gotIDs = %v, want %v", gotIDs, tt.expectedIDs)
			}
			if (gotErr != nil && tt.expectedError == nil) || (gotErr == nil && tt.expectedError != nil) || (gotErr != nil && tt.expectedError != nil && gotErr.Error() != tt.expectedError.Error()) {
				t.Errorf("GetFollowingUserIDs() gotErr = %v, want %v", gotErr, tt.expectedError)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %s", err)
			}
		})
	}
}

func generateSequence(start, end int) []uint {
	var sequence []uint
	for i := start; i <= end; i++ {
		sequence = append(sequence, uint(i))
	}
	return sequence
}

