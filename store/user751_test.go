package store

import (
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/DATA-DOG/go-sqlmock"
)









/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9

FUNCTION_DEF=func NewUserStore(db *gorm.DB) *UserStore 

 */
func BenchmarkNewUserStore(b *testing.B) {
	db := &gorm.DB{}
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
			name: "Valid gorm.DB instance",
			db:   &gorm.DB{},
			want: &UserStore{db: &gorm.DB{}},
		},
		{
			name: "Nil gorm.DB instance",
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

func TestNewUserStorePreservesDBConfig(t *testing.T) {

	db := &gorm.DB{}

	store := NewUserStore(db)

	if !reflect.DeepEqual(store.db, db) {
		t.Error("NewUserStore() did not preserve gorm.DB configuration")
	}
}


/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920

FUNCTION_DEF=func (s *UserStore) Create(m *model.User) error 

 */
func TestUserStoreCreate(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		mockErr error
		wantErr bool
	}{
		{
			name: "Successfully Create a New User",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Create a User with a Duplicate Username",
			user: &model.User{
				Username: "existinguser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockErr: errors.New("Error 1062: Duplicate entry 'existinguser' for key 'username'"),
			wantErr: true,
		},
		{
			name: "Create User with Minimum Required Fields",
			user: &model.User{
				Username: "minuser",
				Email:    "min@example.com",
				Password: "password123",
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "Attempt to Create User with Invalid Email Format",
			user: &model.User{
				Username: "invalidemail",
				Email:    "invalid-email",
				Password: "password123",
			},
			mockErr: errors.New("validation failed: Email is not a valid email address"),
			wantErr: true,
		},
		{
			name: "Database Connection Failure During User Creation",
			user: &model.User{
				Username: "connectionfail",
				Email:    "connection@example.com",
				Password: "password123",
			},
			mockErr: errors.New("failed to connect to database"),
			wantErr: true,
		},
		{
			name: "Create User with Maximum Length Values",
			user: &model.User{
				Username: "maxlengthuser",
				Email:    "maxlength@example.com",
				Password: "verylongpasswordthatreachesmaximumlength",
				Bio:      "This is a very long bio that reaches the maximum allowed length for the bio field in the database",
				Image:    "https://example.com/very/long/image/url/that/reaches/maximum/length/allowed/for/image/field/in/database",
			},
			mockErr: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &gorm.DB{}
			mockDB.Error = tt.mockErr

			us := &UserStore{db: mockDB}

			err := us.Create(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Create() error = %v, wantErr %v", err, tt.wantErr)
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
	tests := []struct {
		name    string
		setup   func(*gorm.DB) (*model.User, *model.User)
		wantErr bool
	}{
		{
			name: "Successful Follow Operation",
			setup: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				userB := &model.User{Username: "userB", Email: "userB@example.com"}
				db.Create(userA)
				db.Create(userB)
				return userA, userB
			},
			wantErr: false,
		},
		{
			name: "Follow User Already Being Followed",
			setup: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				userB := &model.User{Username: "userB", Email: "userB@example.com"}
				db.Create(userA)
				db.Create(userB)
				db.Model(userA).Association("Follows").Append(userB)
				return userA, userB
			},
			wantErr: false,
		},
		{
			name: "Follow Non-Existent User",
			setup: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				userB := &model.User{Username: "userB", Email: "userB@example.com"}
				db.Create(userA)
				return userA, userB
			},
			wantErr: true,
		},
		{
			name: "Follow Self",
			setup: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				db.Create(userA)
				return userA, userA
			},
			wantErr: false,
		},
		{
			name: "Follow with Nil User",
			setup: func(db *gorm.DB) (*model.User, *model.User) {
				userA := &model.User{Username: "userA", Email: "userA@example.com"}
				db.Create(userA)
				return userA, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("failed to open database: %v", err)
			}
			defer db.Close()

			db.AutoMigrate(&model.User{})

			s := &UserStore{db: db}

			userA, userB := tt.setup(db)

			err = s.Follow(userA, userB)

			if (err != nil) != tt.wantErr {
				t.Errorf("Follow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && userB != nil {
				count := db.Model(userA).Association("Follows").Count()
				if count != 1 {
					t.Errorf("Expected 1 follow, got %d", count)
				}

				var followedUser model.User
				db.Model(userA).Association("Follows").Find(&followedUser)
				if followedUser.ID != userB.ID {
					t.Errorf("Expected to follow user with ID %d, but followed user with ID %d", userB.ID, followedUser.ID)
				}
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
			name: "Large number of followed users",
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
			name: "Handling of deleted users",
			user: &model.User{Model: gorm.Model{ID: 1}},
			mockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(2).
					AddRow(3)
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    []uint{2, 3},
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

