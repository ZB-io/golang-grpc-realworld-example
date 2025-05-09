package store

import (
	errors "errors"
	debug "runtime/debug"
	testing "testing"

	"github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

/*
ROOST_METHOD_HASH=UserStore_Create_9495ddb29d
ROOST_METHOD_SIG_HASH=UserStore_Create_18451817fe

FUNCTION_DEF=func (s *UserStore) Create(m *model.User) error // Create create a user
*/
func TestUserStoreCreate(t *testing.T) {
	tests := []struct {
		name         string
		user         *model.User
		mockBehavior func(mock sqlmock.Sqlmock)
		wantErr      bool
	}{
		{
			name: "Successful User Creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "securepassword",
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Database Error During User Creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "securepassword",
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("database error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Invalid User Model",
			user: &model.User{
				Username: "",
				Email:    "",
				Password: "",
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {},
			wantErr:      true,
		},
		{
			name:         "Empty User Model",
			user:         &model.User{},
			mockBehavior: func(mock sqlmock.Sqlmock) {},
			wantErr:      true,
		},
		{
			name: "User Creation with Existing Email",
			user: &model.User{
				Username: "testuser",
				Email:    "existing@example.com",
				Password: "securepassword",
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO \"users\"").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("duplicate email"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "User Creation with Invalid Email Format",
			user: &model.User{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "securepassword",
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {},
			wantErr:      true,
		},
		{
			name: "User Creation with Long Password",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "thisisaverylongpasswordthatisevenlongerthanthemaximumallowedlength",
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {},
			wantErr:      true,
		},
		{
			name: "User Creation with Special Characters in Username",
			user: &model.User{
				Username: "test@user",
				Email:    "test@example.com",
				Password: "securepassword",
			},
			mockBehavior: func(mock sqlmock.Sqlmock) {},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("sqlmock", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
			}
			defer gormDB.Close()

			tt.mockBehavior(mock)

			userStore := &UserStore{db: gormDB}

			err = userStore.Create(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.Create() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("Test passed for scenario: %s", tt.name)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
